from dataclient.result import Result, deserialize_result
import grpc
import uuid
import capnp
import dataclient.data_capnp
import pytz
import pandas as pd
import dataclient.mdal_pb2 as mdal_pb2
import dataclient.mdal_pb2_grpc as mdal_pb2_grpc

agg_funcs = {
    "RAW": mdal_pb2.RAW,
    "MEAN": mdal_pb2.MEAN,
    "MIN": mdal_pb2.MIN,
    "MAX": mdal_pb2.MAX,
    "COUNT": mdal_pb2.COUNT,
    "SUM": mdal_pb2.SUM,
}
def parse_agg_func(name):
    return mdal_pb2.AggFunc.Value(name.upper())

class MDALClient:
    def __init__(self, caddr):
        self.channel = grpc.insecure_channel(caddr)
        self.stub = mdal_pb2_grpc.MDALStub(self.channel)

    def query(self, request):
        aggs = {}
        for varname, aggfunclist in request["Aggregation"].items():
            aggs[varname] = mdal_pb2.Aggregation(funcs=[parse_agg_func(func) for func in aggfunclist])
        vardefs = {}
        for varname, defn in request["Variables"].items():
            vardefs[varname] = mdal_pb2.Variable(
                name = varname,
                definition = defn["Definition"],
                units = defn.get("Units",None),
            )
        params = mdal_pb2.DataQueryRequest(
            composition = request["Composition"],
            aggregation = aggs,
            variables = vardefs,
            time = mdal_pb2.TimeParams(
                start = request["Time"]["Start"],
                end = request["Time"]["End"],
                window = request["Time"]["Window"],
                aligned = request["Time"]["Aligned"],
            ),
        )
        tz = pytz.timezone("US/Pacific")
        responses = self.stub.DataQuery(params, timeout=120)
        for resp in responses:
            if resp.error != "":
                raise Exception(resp.error)
            result = Result(resp)
            return result
        #    uuids = [str(uuid.UUID(bytes=x)) for x in resp.uuids]
        #    data = data_capnp.StreamCollection.from_bytes_packed(resp.arrow)
        #    if hasattr(data, 'times') and len(data.times):
        #       times = list(data.times)
        #       if len(times) == 0:
        #           return pd.DataFrame(columns=uuids)
        #       df = pd.DataFrame(index=pd.to_datetime(times, unit='ns', utc=False))
        #       for idx, s in enumerate(data.streams):
        #           df[uuids[idx]] = s.values
        #       df.index = df.index.tz_localize(pytz.utc).tz_convert(tz)
        #       return df
        #    else:
        #       df = pd.DataFrame()
        #       for idx, s in enumerate(data.streams):
        #           if hasattr(s, 'times'):
        #               newdf = pd.DataFrame(list(s.values), index=list(s.times), columns=[uuids[idx]])
        #               newdf.index = pd.to_datetime(newdf.index, unit='ns').tz_localize(pytz.utc).tz_convert(tz)
        #               df = df.join(newdf, how='outer')
        #           else:
        #               raise Exception("Does this ever happen? Tell gabe!")
        #       return df

if __name__ == '__main__':
    m = MDALClient("corbusier.cs.berkeley.edu:8088")

    for windowsize in ["96h","24h","12h","6h","3h","1h","30m","15m","10m","5m","1m","30s"]:
    #for windowsize in ["1h","30m","15m","10m","5m","1m","30s"]:
        request = {
            "Composition": ["temp"],
            "Aggregation": {
                "temp": ["MEAN"],
            },
            "Variables": {
                "temp": {
                    "Definition": """SELECT ?temp ?temp_uuid FROM * WHERE {
                        ?temp rdf:type/rdfs:subClassOf* brick:Electric_Meter .
                        ?temp bf:uuid ?temp_uuid
                    };""",
                },
            },
            "Time": {
                "Start":  "2017-01-01T10:00:00-07:00",
                "End":  "2018-05-12T10:00:00-07:00",
                "Window": windowsize,
                "Aligned": True,
            },
        }
        try:
            resp = m.query(request)
        except Exception as e:
            print('ERR', e)
            continue
        print(resp)
