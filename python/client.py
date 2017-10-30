import msgpack
import time
import threading
import random
from bw2python import ponames
from bw2python.bwtypes import PayloadObject
RAW  = 0
MEAN = 1 << 0
MIN = 1 <<  1
MAX = 1 << 2
COUNT = 1 << 3

query = {
    "Composition": ["meter", "temp"],
    "Selectors": [RAW, RAW],
    "Variables": [
        {"Name": "meter",
         "Definition": "SELECT ?meter_uuid WHERE { ?meter rdf:type/rdfs:subClassOf* brick:Meter . ?meter bf:uuid ?meter_uuid . };",
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
        "WindowSize": 0,
        "Aligned": False,
    },
    "Params": {
        "Statistical": False,
        "Window": False,
    },
}

if __name__ == '__main__':
    from xbos import get_client
    c = get_client()

    DEFAULT_TIMEOUT=10
    def do_query(timeout=DEFAULT_TIMEOUT, values_only=True):
        #nonce = str(random.randint(0, 2**32))
        ev = threading.Event()
        response = {}
        def _handleresult(msg):
            print [po.type_dotted for po in msg.payload_objects]
            for po in msg.payload_objects:
                if po.type_dotted == (2,0,10,4):
                    data = msgpack.unpackb(po.content)
                    print data

            ev.set()
        c.subscribe("{0}/s.mdal/_/i.mdal/signal/result".format("scratch.ns"), _handleresult)
        po = PayloadObject((2,0,10,3), None, msgpack.packb(query))
        c.publish("{0}/s.mdal/_/i.mdal/slot/query".format("scratch.ns"), payload_objects=(po,))
        ev.wait(timeout)
        return response
    do_query()
