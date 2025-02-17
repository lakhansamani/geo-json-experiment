// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.6.1
//   protoc               v5.27.3
// source: apis/layer/layer.proto

/* eslint-disable */
import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import { ListValue, Struct } from "../../google/protobuf/struct";

export const protobufPackage = "layer.v1";

/** Message for GeoJSON Geometry */
export interface Geometry {
  type: string;
  /** Coordinates stored as nested arrays using google.protobuf.ListValue */
  coordinates: Array<any> | undefined;
}

/** GeoJSON Feature */
export interface Feature {
  type: string;
  bbox: number[];
  geometry: Geometry | undefined;
  properties: { [key: string]: any } | undefined;
}

function createBaseGeometry(): Geometry {
  return { type: "", coordinates: undefined };
}

export const Geometry: MessageFns<Geometry> = {
  encode(message: Geometry, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.coordinates !== undefined) {
      ListValue.encode(ListValue.wrap(message.coordinates), writer.uint32(18).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): Geometry {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGeometry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.type = reader.string();
          continue;
        }
        case 2: {
          if (tag !== 18) {
            break;
          }

          message.coordinates = ListValue.unwrap(ListValue.decode(reader, reader.uint32()));
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Geometry {
    return {
      type: isSet(object.type) ? globalThis.String(object.type) : "",
      coordinates: globalThis.Array.isArray(object.coordinates) ? [...object.coordinates] : undefined,
    };
  },

  toJSON(message: Geometry): unknown {
    const obj: any = {};
    if (message.type !== "") {
      obj.type = message.type;
    }
    if (message.coordinates !== undefined) {
      obj.coordinates = message.coordinates;
    }
    return obj;
  },

  create(base?: DeepPartial<Geometry>): Geometry {
    return Geometry.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<Geometry>): Geometry {
    const message = createBaseGeometry();
    message.type = object.type ?? "";
    message.coordinates = object.coordinates ?? undefined;
    return message;
  },
};

function createBaseFeature(): Feature {
  return { type: "", bbox: [], geometry: undefined, properties: undefined };
}

export const Feature: MessageFns<Feature> = {
  encode(message: Feature, writer: BinaryWriter = new BinaryWriter()): BinaryWriter {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    writer.uint32(18).fork();
    for (const v of message.bbox) {
      writer.float(v);
    }
    writer.join();
    if (message.geometry !== undefined) {
      Geometry.encode(message.geometry, writer.uint32(26).fork()).join();
    }
    if (message.properties !== undefined) {
      Struct.encode(Struct.wrap(message.properties), writer.uint32(34).fork()).join();
    }
    return writer;
  },

  decode(input: BinaryReader | Uint8Array, length?: number): Feature {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFeature();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1: {
          if (tag !== 10) {
            break;
          }

          message.type = reader.string();
          continue;
        }
        case 2: {
          if (tag === 21) {
            message.bbox.push(reader.float());

            continue;
          }

          if (tag === 18) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.bbox.push(reader.float());
            }

            continue;
          }

          break;
        }
        case 3: {
          if (tag !== 26) {
            break;
          }

          message.geometry = Geometry.decode(reader, reader.uint32());
          continue;
        }
        case 4: {
          if (tag !== 34) {
            break;
          }

          message.properties = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        }
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skip(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Feature {
    return {
      type: isSet(object.type) ? globalThis.String(object.type) : "",
      bbox: globalThis.Array.isArray(object?.bbox) ? object.bbox.map((e: any) => globalThis.Number(e)) : [],
      geometry: isSet(object.geometry) ? Geometry.fromJSON(object.geometry) : undefined,
      properties: isObject(object.properties) ? object.properties : undefined,
    };
  },

  toJSON(message: Feature): unknown {
    const obj: any = {};
    if (message.type !== "") {
      obj.type = message.type;
    }
    if (message.bbox?.length) {
      obj.bbox = message.bbox;
    }
    if (message.geometry !== undefined) {
      obj.geometry = Geometry.toJSON(message.geometry);
    }
    if (message.properties !== undefined) {
      obj.properties = message.properties;
    }
    return obj;
  },

  create(base?: DeepPartial<Feature>): Feature {
    return Feature.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<Feature>): Feature {
    const message = createBaseFeature();
    message.type = object.type ?? "";
    message.bbox = object.bbox?.map((e) => e) || [];
    message.geometry = (object.geometry !== undefined && object.geometry !== null)
      ? Geometry.fromPartial(object.geometry)
      : undefined;
    message.properties = object.properties ?? undefined;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export interface MessageFns<T> {
  encode(message: T, writer?: BinaryWriter): BinaryWriter;
  decode(input: BinaryReader | Uint8Array, length?: number): T;
  fromJSON(object: any): T;
  toJSON(message: T): unknown;
  create(base?: DeepPartial<T>): T;
  fromPartial(object: DeepPartial<T>): T;
}
