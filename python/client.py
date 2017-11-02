import msgpack
import pandas as pd
import uuid
import time
from datetime import timedelta, datetime
import threading
import random
import capnp
import data_capnp
from bw2python import ponames
from bw2python.bwtypes import PayloadObject
from bw2python.client import Client
from xbos.util import pretty_print_timedelta, TimeoutException
RAW  = 0
MEAN = 1 << 1
MIN = 1 <<  2
MAX = 1 << 3
COUNT = 1 << 4

DEFAULT_TIMEOUT=30
class MDALClient(object):
    def __init__(self, url, client=None, timeout=30):
        if client is None:
            client = Client()
            client.vk = client.setEntityFromEnviron()
            client.overrideAutoChainTo(True)
        if not isinstance(client, Client):
            raise TypeError("first argument must be bw2python.client.Client or None")
        self.c = client
        self.vk = client.vk
        self.url = url

        # check liveness
        responses = self.c.query("{0}/*/s.mdal/!meta/lastalive".format(url))
        for resp in responses:
            # get the metadata records from the response
            md_records = filter(lambda po: po.type_dotted == ponames.PODFSMetadata, resp.payload_objects)
            # get timestamp field from the first metadata record
            last_seen_timestamp = msgpack.unpackb(md_records[0].content).get('ts')
            # get how long ago that was
            now = time.time()*1e9 # nanoseconds
            # convert to microseconds and get the timedelta
            diff = timedelta(microseconds = (now - last_seen_timestamp)/1e3)
            print "Saw [{0}] MDAL {1}".format(self.url, pretty_print_timedelta(diff))
            if diff.total_seconds() > timeout:
                raise TimeoutException("MDAL at {0} is too old".format(self.url))
    def do_query(self, query, timeout=DEFAULT_TIMEOUT):
        #nonce = str(random.randint(0, 2**32))
        ev = threading.Event()
        response = {}
        def _handleresult(msg):
            got_response = False
            for po in msg.payload_objects:
                if po.type_dotted == (2,0,10,4):
                    data = msgpack.unpackb(po.content)
                    uuids = [uuid.UUID(bytes=x) for x in data['Rows']]
                    data = data_capnp.StreamCollection.from_bytes_packed(data['Data'])
                    if hasattr(data, 'times'):
                        times = list(data.times)
                        df = pd.DataFrame(index=pd.to_datetime(times, unit='ns'))
                        for idx, s in enumerate(data.streams):
                            df[uuids[idx]] = s.values
                        response['df'] = df
                        got_response = True
                    else:
                        for idx, s in enumerate(data.streams):
                            if hasattr(s, 'times'):
                                s = pd.Series(s.values, pd.to_datetime(s.times, unit='ns'))
                            else:
                                s = pd.Series(s.values)
                            response[uuids[idx]] = s
                        got_response = True
            if got_response:
                ev.set()
        h = self.c.subscribe("{0}/s.mdal/_/i.mdal/signal/{1}".format("scratch.ns", self.vk[:-1]), _handleresult)
        po = PayloadObject((2,0,10,3), None, msgpack.packb(query))
        self.c.publish("{0}/s.mdal/_/i.mdal/slot/query".format("scratch.ns"), payload_objects=(po,))
        ev.wait(timeout)
        self.c.unsubscribe(h)
        return response

if __name__ == '__main__':
    query = {
        "Composition": ["meter", "temp"],
        "Selectors": [MEAN, MEAN],
        "Variables": [
            {"Name": "meter",
             "Definition": "SELECT ?meter_uuid WHERE { ?meter rdf:type/rdfs:subClassOf* brick:Electric_Meter . ?meter bf:uuid ?meter_uuid . };",
             "Units": "kW",
            },
            {"Name": "temp",
             "Definition": "SELECT ?temp_uuid WHERE { ?temp rdf:type/rdfs:subClassOf* brick:Temperature_Sensor . ?temp bf:uuid ?temp_uuid . };",
             "Units": "F",
            },
        ],
        "Time": {
            "T0": "2017-08-01 00:00:00",
            "T1": "2017-08-08 00:00:00",
            "WindowSize": int(1e9*60*5),
            "Aligned": True,
        },
        "Params": {
            "Statistical": False,
            "Window": True,
        },
    }

    c = MDALClient("scratch.ns")
    df = c.do_query(query)['df']
    print df.describe()
