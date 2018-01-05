# Metadata-driven Data Access Layer

Returns resampled/aligned timeseries data from Brick queries + UUIDs

Requests look like:

```python
query = {
    "Composition": ["temp"],
    "Selectors": [MEAN],
    "Variables": [
        {"Name": "meter",
         "Definition": "SELECT ?meter_uuid WHERE { ?meter rdf:type/rdfs:subClassOf* brick:Electric_Meter . ?meter bf:uuid ?meter_uuid . };",
         "Units": "kW",
        },
        {"Name": "temp",
         "Definition": "SELECT ?temp_uuid WHERE { ?temp rdf:type/rdfs:subClassOf* brick:Temperature_Sensor . ?temp bf:uuid ?temp_uuid . };",
         "Units": "C",
        },
    ],
    "Time": {
        "T0": "2017-07-21 00:00:00",
        "T1": "2017-08-30 00:00:00",
        "WindowSize": '2h',
        "Aligned": True,
    },
}
```

- `Composition`: the order of variables and UUIDs to be included in the response matrix. Variables are defined in `Variables` key (below) and resolve to a set of UUIDs. UUIDs are the pointers used by the timeseries database and represent a single timeseries sequence
- `Selectors`: for each timeseries stream, we can get the raw data, or we can get resampled min, mean and/or max (as well as bucket count). Each item in the `Composition` list has a corresponding selector. This is a set of flags:
	- `MEAN`: selects the mean stream
	- `MAX`: selects the max stream
	- `MIN`: selects the min stream
	- `COUNT`: selects the count stream (how many readings in each resampled bucket)
	- `RAW`: selects the raw data. This cannot be resampled and is mutually exclusive with
	  the above flags
	Combine flags like `MEAN|MAX`
- `Variables`: each variable mentioned in `Composition` has a definition here. Each variable needs the following fields:
    - `Name`: the name of the variable. Refer to the variable in `Composition` using this name
    - `Definition`: the BRICK query that will be resolved to a set of UUIDs. The returned variables need to end in `_uuid`
    - `Units`: the desired units for the stream; MDAL will perform the unit conversion if possible
- `Time`: the temporal parameters of the data query
    - `T0`: the first half of the range (inclusive)
    - `T1`: the second half of the range (inclusive). These will be reordered appropriately by MDAL
    - `WindowSize`: window size as a Go-style duration, e.g. "5m", "1h", "3d". 
        - Supported units are: `d`,`h`,`m`,`s`,`us`,`ms`,`ns`
    - `Aligned`: if true, then all timesetamps will be the same, else each stream (UUID) will have its own timestamps. Its best to leave this as `True`
- `Params`: don't touch this for now

Supported unit conversions (case insensitive):
- `w/watt`, `kw/kilowatt`
- `wh/watthour`, `kwh`,`kilowatthour`
- `c/celsius`,`f/fahrenheit`,`k/kelvin`


## Python Client

```python
from xbos.services import mdal
client = mdal.BOSSWAVEMDALClient("scratch.ns")
query = {
    "Composition": ["temp"],
    "Selectors": [mdal.MEAN],
    "Variables": [
        {"Name": "temp",
         "Definition": "SELECT ?temp_uuid WHERE { ?temp rdf:type/rdfs:subClassOf* brick:Temperature_Sensor . ?temp bf:uuid ?temp_uuid . };",
         "Units": "C",
        },
    ],
    "Time": {
        "T0": "2017-07-21 00:00:00",
        "T1": "2017-08-21 00:00:00",
        "WindowSize": '30m',
        "Aligned": True,
    },
    "Params": {
        "Statistical": False,
        "Window": True,
    },
}
resp = client.do_query(query,timeout=300)
df = resp['df']
print df.describe()
```
