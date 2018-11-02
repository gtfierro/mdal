import uuid
import pytz
import pandas as pd
import capnp
import dataclient.data_capnp as data_capnp
from dataclient.cached_property import cached_property
from dataclient.mdal_pb2 import DataQueryResponse

def deserialize_result(s):
    return Result(DataQueryResponse.FromString(s))

class Result(object):
    def __init__(self, response_object, tz = pytz.timezone("US/Pacific")):
        self.tz = tz
        self.serialized = response_object.SerializeToString()
        self.mapping = {}
        self.context = {}
        self.uuids = response_object.uuids
        self.data = data_capnp.StreamCollection.from_bytes_packed(response_object.arrow)
        for row in response_object.context:
            self.context[row.uuid] = dict(row.row)
        for k, v in response_object.mapping.items():
            self.mapping[k] = v.uuids #[str(uuid.UUID(bytes=x)) for x in v.uuids]
        self.df = self._build_df()

    def _build_df(self):
        if hasattr(self.data, 'times') and len(self.data.times):
            times = list(self.data.times)
            if len(times) == 0:
                return pd.DataFrame(columns=self.uuids)
            df = pd.DataFrame(index=pd.to_datetime(times, unit='ns', utc=False))
            for idx, s in enumerate(self.data.streams):
                df[self.uuids[idx]] = s.values
            df.index = df.index.tz_localize(pytz.utc).tz_convert(self.tz)
            return df
        else:
            df = pd.DataFrame()
            for idx, s in enumerate(self.data.streams):
                if hasattr(s, 'times'):
                    newdf = pd.DataFrame(list(s.values), index=list(s.times), columns=[self.uuids[idx]])
                    newdf.index = pd.to_datetime(newdf.index, unit='ns').tz_localize(pytz.utc).tz_convert(self.tz)
                    df = df.join(newdf, how='outer')
                else:
                    raise Exception("Does this ever happen? Tell gabe!")
            return df

    def join_on(self, uuid, key, varnames=None):
        uuids = [uuid]
        joinval = self.context.get(uuid)
        if joinval is None:
            raise KeyError("No value for key {0} on uuid {1}".format(key, uuid))
        if varnames is not None:
            for varname in varnames:
                for uuid in self.mapping[varname]:
                    metadata = self.context[uuid]
                    if metadata.get(key) == joinval[key]:
                        uuids.append(uuid)
        else:
            for uuid, metadata in self.context.items():
                if metadata.get(key) == joinval[key]:
                    uuids.append(uuid)
        return uuids

        

    @cached_property()
    def context(self):
        for row in self.response_object.context:
            self._context[row.uuid] = dict(row.row)
        return self._context

    @cached_property()
    def mapping(self):
        for k, v in self.response_object.mapping.items():
            self._mapping[k] = [str(uuid.UUID(bytes=x)) for x in v.uuids]
        return self._mapping
