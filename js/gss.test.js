// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

const { serialize, serializeArray, deserialize, convert, formats } = global.gss;

const testObject = {
  "a": "x",
  "b": "y",
  "c": "z"
};

const testArray = [
  {
    "a": "x",
    "b": "y",
    "c": "z"
  },
  {
    "b": "g",
    "c": "h",
    "d": "i"
  }
];

function log(str) {
  console.log(str.replace(/\n/g, "\\n").replace(/\t/g, "\\t").replace(/"/g, "\\\""));
}

describe('gss', () => {

  it('checks the available formats', () => {
    expect(formats).toEqual(["bson", "csv", "go", "json", "jsonl", "properties", "tags", "toml", "tsv", "hcl", "hcl2", "yaml"]);
  });

});

describe('serialize : object', () => {

  it('serializes an object to csv', () => {
    var { str, err } = serialize(testObject, "csv", {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a,b,c\nx,y,z\n");
  });

  it('serializes an object to go', () => {
    var { str, err } = serialize(testObject, "go");
    expect(err).toBeNull();
    expect(str).toEqual("map[string]interface {}{\"a\":\"x\", \"b\":\"y\", \"c\":\"z\"}");
  });

  it('serializes an object to json', () => {
    var { str, err } = serialize(testObject, "json");
    expect(err).toBeNull();
    expect(str).toEqual("{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"}");
  });

  it('serializes an object to json (pretty)', () => {
    var { str, err } = serialize(testObject, "json", {"pretty": true});
    expect(err).toBeNull();
    expect(str).toEqual("{\n  \"a\": \"x\",\n  \"b\": \"y\",\n  \"c\": \"z\"\n}");
  });

  it('serializes an object to jsonl', () => {
    var { str, err } = serialize(testObject, "jsonl");
    expect(err).toBeNull();
    expect(str).toEqual("{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"}\n");
  });

  it('serializes an object to properties', () => {
    var { str, err } = serialize(testObject, "properties");
    expect(err).toBeNull();
    expect(str).toEqual("a=x\nb=y\nc=z");
  });

  it('serializes an object to tags', () => {
    var { str, err } = serialize(testObject, "tags");
    expect(err).toBeNull();
    expect(str).toEqual("a=x b=y c=z");
  });

  it('serializes an object to toml', () => {
    var { str, err } = serialize(testObject, "toml");
    expect(err).toBeNull();
    expect(str).toEqual("a = \"x\"\nb = \"y\"\nc = \"z\"\n");
  });

  it('serializes an object to tsv', () => {
    var { str, err } = serialize(testObject, "tsv");
    expect(err).toBeNull();
    expect(str).toEqual("a\tb\tc\nx\ty\tz\n");
  });

  it('serializes an object to yaml', () => {
    var { str, err } = serialize(testObject, "yaml");
    expect(err).toBeNull();
    expect(str).toEqual("a: x\nb: \"y\"\nc: z\n");
  });

});

describe('serialize : array', () => {

  it('serializes an array to csv', () => {
    var { str, err } = serialize(testArray, "csv", {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a,b,c\nx,y,z\n,g,h\n");
  });

  it('serializes an array to csv (reversed)', () => {
    var { str, err } = serialize(testArray, "csv", {"sorted": true, "reversed": true});
    expect(err).toBeNull();
    expect(str).toEqual("c,b,a\nz,y,x\nh,g,\n");
  });

  it('serializes an array to csv (expand header)', () => {
    var { str, err } = serialize(testArray, "csv", {"sorted": true, "expandHeader": true});
    expect(err).toBeNull();
    expect(str).toEqual("a,b,c,d\nx,y,z,\n,g,h,i\n");
  });

  it('serializes an array to go', () => {
    var { str, err } = serialize(testArray, "go");
    expect(err).toBeNull();
    expect(str).not.toBeNull(); // hard to test, since not sorted.
  });

  it('serializes an array to json', () => {
    var { str, err } = serialize(testArray, "json");
    expect(err).toBeNull();
    expect(str).toEqual("[{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"},{\"b\":\"g\",\"c\":\"h\",\"d\":\"i\"}]");
  });

  it('serializes an array to json (pretty)', () => {
    var { str, err } = serialize(testArray, "json", {"pretty": true});
    //console.error(str.replace(/\n/g, "\\n").replace(/"/g, "\\\""));
    expect(err).toBeNull();
    expect(str).toEqual("[\n  {\n    \"a\": \"x\",\n    \"b\": \"y\",\n    \"c\": \"z\"\n  },\n  {\n    \"b\": \"g\",\n    \"c\": \"h\",\n    \"d\": \"i\"\n  }\n]");
  });

  it('serializes an array to jsonl', () => {
    var { str, err } = serialize(testArray, "jsonl");
    expect(err).toBeNull();
    expect(str).toEqual("{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"}\n{\"b\":\"g\",\"c\":\"h\",\"d\":\"i\"}\n");
  });

  it('serializes an array to properties (error)', () => {
    var { str, err } = serialize(testArray, "properties", {"sorted": true});
    expect(err).toEqual("error serializing input object: error serializing: error writing properties: type \"[]interface {}\" is of invalid kind, expecting one of [\"map\" \"struct\"]");
    expect(str).toBeNull();
  });

  it('serializes an array to tags', () => {
    var { str, err } = serialize(testArray, "tags", {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a=x b=y c=z\nb=g c=h d=i\n");
  });

  it('serializes an array to toml (error)', () => {
    var { str, err } = serialize(testArray, "toml");
    expect(err).toEqual("error serializing input object: error serializing: error marshaling TOML bytes: toml: top-level values must be Go maps or structs");
    expect(str).toBeNull();
  });

  it('serializes an array to tsv', () => {
    var { str, err } = serialize(testArray, "tsv");
    expect(err).toBeNull();
    expect(str).toEqual("a\tb\tc\nx\ty\tz\n\tg\th\n");
  });

  it('serializes an array to yaml', () => {
    var { str, err } = serialize(testArray, "yaml");
    expect(err).toBeNull();
    expect(str).toEqual("- a: x\n  b: \"y\"\n  c: z\n- b: g\n  c: h\n  d: i\n");
  });

});

describe('deserialize : object', () => {

  it('deserializes an object from csv', () => {
    var { obj, err } = deserialize("a,b,c\nx,y,z\n", "csv", {});
    expect(err).toBeNull();
    expect(obj).toEqual([{"a": "x", "b": "y", "c": "z"}]);
  });

  /*

  it('serializes an object from go', () => {
    var { str, err } = serialize(testObject, "go");
    expect(err).toBeNull();
    expect(str).toEqual("map[string]interface {}{\"a\":\"x\", \"b\":\"y\", \"c\":\"z\"}");
  });

  */

  it('deserializes an object from json', () => {
    var { obj, err } = deserialize("{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"}", "json");
    expect(err).toBeNull();
    expect(obj).toEqual(testObject);
  });

  it('deserializes an object from jsonl', () => {
    var { obj, err } = deserialize("{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"}", "jsonl");
    expect(err).toBeNull();
    expect(obj).toEqual([testObject]);
  });

  it('deserializes an object from properties', () => {
    var { obj, err } = deserialize("a=x\nb=y\nc=z", "properties");
    expect(err).toBeNull();
    expect(obj).toEqual(testObject);
  });

  it('deserializes an object from tags', () => {
    var { obj, err } = deserialize("a=x b=y c=z", "tags");
    expect(err).toBeNull();
    expect(obj).toEqual([testObject]);
  });

  it('deserializes an object from toml', () => {
    var { obj, err } = deserialize("a = \"x\"\nb = \"y\"\nc = \"z\"\n", "toml");
    expect(err).toBeNull();
    expect(obj).toEqual(testObject);
  });

  it('deserializes an object from tsv', () => {
    var { obj, err } = deserialize("a\tb\tc\nx\ty\tz\n", "tsv");
    expect(err).toBeNull();
    expect(obj).toEqual([testObject]);
  });

  it('deserializes an object to yaml', () => {
    var { obj, err } = deserialize("a: x\nb: \"y\"\nc: z\n", "yaml");
    expect(err).toBeNull();
    expect(obj).toEqual(testObject);
  });

});

describe('convert : object', () => {

  it('convert an object from csv to json', () => {
    var { str, err } = convert("a,b,c\nx,y,z\n", "csv", "json", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("[{\"a\":\"x\",\"b\":\"y\",\"c\":\"z\"}]");
  });

  it('convert an object from csv to tsv', () => {
    var { str, err } = convert("a,b,c\nx,y,z\n", "csv", "tsv", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a\tb\tc\nx\ty\tz\n");
  });

  it('convert an object from csv to yaml', () => {
    var { str, err } = convert("a,b,c\nx,y,z\n", "csv", "yaml", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("- a: x\n  b: \"y\"\n  c: z\n");
  });

  it('convert an object from json to bson and back', () => {
    var { str, err } = convert(JSON.stringify(testObject), "json", "bson");
    expect(err).toBeNull();
    // TODO: hex or base64 encode string so we can test
    //expect(Buffer.to(str, 'hex')).toEqual("a=x\nb=y\nc=z");
    var { str, err } = convert(str, "bson", "json");
    expect(err).toBeNull();
    expect(str).toEqual(JSON.stringify(testObject));
  });

  it('convert an object from json to properties and back', () => {
    var { str, err } = convert(JSON.stringify(testObject), "json", "properties", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a=x\nb=y\nc=z");
    var { str, err } = convert(str, "properties", "json", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual(JSON.stringify(testObject));
  });

  it('convert an object from json to tags and back', () => {
    var { str, err } = convert(JSON.stringify(testObject), "json", "tags", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a=x b=y c=z");
    var { str, err } = convert(str, "tags", "json", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual(JSON.stringify([testObject])); // to support streaming, deserializing tags returns an array
  });

  it('convert an object from json to toml and back', () => {
    var { str, err } = convert(JSON.stringify(testObject), "json", "toml", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a = \"x\"\nb = \"y\"\nc = \"z\"\n");
    var { str, err } = convert(str, "toml", "json", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual(JSON.stringify(testObject));
  });

  it('convert an object from json to yaml and back', () => {
    var { str, err } = convert(JSON.stringify(testObject), "json", "yaml", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual("a: x\nb: \"y\"\nc: z\n");
    var { str, err } = convert(str, "yaml", "json", undefined, {"sorted": true});
    expect(err).toBeNull();
    expect(str).toEqual(JSON.stringify(testObject));
  });

});
