# Python Client

Installation:

```
pip install dataclient
```

Usage:

```python
import dataclient
m = dataclient.MDALClient("corbusier.cs.berkeley.edu:8088")

request = {
    "Composition": ["temp"],
    "Aggregation": {
        "meter": ["MEAN"],
        "temp": ["MEAN"],
    },
    "Variables": {
        "meter": {
            "Definition": """SELECT ?meter ?meter_uuid WHERE {
                ?meter rdf:type brick:Building_Electric_Meter .
                ?meter bf:uuid ?meter_uuid
            };""" ,
        },
        "temp": {
            "Definition": """SELECT ?t ?t_uuid FROM ciee WHERE {
                ?t rdf:type/rdfs:subClassOf* brick:Temperature_Sensor .
                ?t bf:uuid ?t_uuid
            };""" ,
        },
    },
    "Time": {
        "Start":  "2018-01-01T10:00:00-07:00",
        "End":  "2018-08-12T10:00:00-07:00",
        "Window": '5min',
        "Aligned": True,
    },
}

resp = m.query(request)
print(resp.df)
```

Fields:
- `Composition`: the order of variables and UUIDs to be included in the response matrix. Variables are defined in `Variables` key (below) and resolve to a set of UUIDs. UUIDs are the pointers used by the timeseries database and represent a single timeseries sequence. If you want to apply unit conversion to the UUIDs, you will need to have instead define the UUIDs using a variable definition below.
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
    - `Definition`: the Brick query that will be resolved to a set of UUIDs. The returned variables need to end in `_uuid`
    - `UUIDS`: a list of additional UUIDs to append to this variable definition (can be used in addition or instead of the Brick query above). To perform unit conversion on data for these UUIDs, put those UUIDs here.
    - `Units`: the desired units for the stream; MDAL will perform the unit conversion if possible
- `Time`: the temporal parameters of the data query
    - `T0`: the first half of the range (inclusive)
    - `T1`: the second half of the range (inclusive). These will be reordered appropriately by MDAL
    - `WindowSize`: window size as a Go-style duration, e.g. "5m", "1h", "3d".
        - Supported units are: `d`,`h`,`m`,`s`,`us`,`ms`,`ns`
    - `Aligned`: if true, then all timesetamps will be the same, else each stream (UUID) will have its own timestamps. Its best to leave this as `True`

Supported unit conversions (case insensitive):
- `w/watt`, `kw/kilowatt`
- `wh/watthour`, `kwh`,`kilowatthour`
- `c/celsius`,`f/fahrenheit`,`k/kelvin`
