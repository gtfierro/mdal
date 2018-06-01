import grpc
import mdal_pb2
import mdal_pb2_grpc

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
        print aggs
        vardefs = {}
        for varname, defn in request["Variables"].items():
            vardefs[varname] = mdal_pb2.Variable(
                definition = defn["Definition"],
                units = defn["Units"],
            )
        print vardefs
        params = mdal_pb2.DataQueryRequest(
            composition = request["Composition"],
            aggregation = aggs,
            variables = vardefs,
            time = mdal_pb2.TimeParams(
                start = "2018-05-02T10:00:00-07:00",
                end = "2018-05-01T10:00:00-07:00",
                window = "30m",
                aligned = True,
            ),
        )
        print params
        resp = self.stub.DataQuery(params)
        print resp


if __name__ == '__main__':
    m = MDALClient("localhost:8088")
    request = {
        "Composition": ["temp"],
        "Aggregation": {
            "temp": ["MEAN"],
        },
        "Variables": {
            "temp": {
                "Definition": """SELECT ?temp ?temp_uuid FROM ciee WHERE {
                    ?temp rdf:type/rdfs:subClassOf* brick:Temperature_Sensor .
                    ?temp bf:uuid ?temp_uuid
                };""",
                "Units": "F",
            },
        },
    }
    m.query(request)
