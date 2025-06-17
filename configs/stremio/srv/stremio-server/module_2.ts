// Continuation of Moment.js and EBML related code

// Export the EBML ID mapping
export const byEbmlID: Record<number, any> = {
    128: {
        name: "ChapterDisplay",
        level: 4,
        type: "m",
        multiple: true,
        minver: 1,
        webm: true,
        description: "Contains all possible strings to use for the chapter display."
    },
    131: {
        name: "TrackType",
        level: 3,
        type: "u",
        mandatory: true,
        minver: 1,
        range: "1-254",
        description: "A set of track types coded on 8 bits (1: video, 2: audio, 3: complex, 0x10: logo, 0x11: subtitle, 0x12: buttons, 0x20: control)."
    },
    133: {
        name: "ChapString",
        cppname: "ChapterString",
        level: 5,
        type: "8",
        mandatory: true,
        minver: 1,
        webm: true,
        description: "Contains the string to use as the chapter atom."
    },
    134: {
        name: "CodecID",
        level: 3,
        type: "s",
        mandatory: true,
        minver: 1,
        description: "An ID corresponding to the codec, see the codec page for more info."
    },
    136: {
        name: "FlagDefault",
        cppname: "TrackFlagDefault",
        level: 3,
        type: "u",
        mandatory: true,
        minver: 1,
        default: 1,
        range: "0-1",
        description: "Set if that track (audio, video or subs) SHOULD be active if no language found matches the user preference. (1 bit)"
    },
    137: {
        name: "ChapterTrackNumber",
        level: 5,
        type: "u",
        mandatory: true,
        multiple: true,
        minver: 1,
        webm: false,
        range: "not 0",
        description: "UID of the Track to apply this chapter too. In the absense of a control track, choosing this chapter will select the listed Tracks and deselect unlisted tracks. Absense of this element indicates that the Chapter should be applied to any currently used Tracks."
    },
    145: {
        name: "ChapterTimeStart",
        level: 4,
        type: "u",
        mandatory: true,
        minver: 1,
        webm: true,
        description: "Timestamp of the start of Chapter (not scaled)."
    },
    146: {
        name: "ChapterTimeEnd",
        level: 4,
        type: "u",
        minver: 1,
        webm: false,
        description: "Timestamp of the end of Chapter (timestamp excluded, not scaled)."
    },
    150: {
        name: "CueRefTime",
        level: 5,
        type: "u",
        mandatory: true,
        minver: 2,
        webm: false,
        description: "Timestamp of the referenced Block."
    },
    151: {
        name: "CueRefCluster",
        level: 5,
        type: "u",
        mandatory: true,
        webm: false,
        description: "The Position of the Cluster containing the referenced Block."
    },
    152: {
        name: "ChapterFlagHidden",
        level: 4,
        type: "u",
        mandatory: true,
        minver: 1,
        webm: false,
        default: 0,
        range: "0-1",
        description: "If a chapter is hidden (1), it should not be available to the user interface (but still to Control Tracks; see flag notes). (1 bit)"
    },
    // ... continued with more EBML mappings
    16980: {
        name: "ContentCompAlgo",
        level: 6,
        type: "u",
        mandatory: true,
        minver: 1,
        webm: false,
        default: 0,
        description: "The compression algorithm used. Algorithms that have been specified so far are: 0 - zlib, 3 - Header Stripping"
    },
    16981: {
        name: "ContentCompSettings",
        level: 6,
        type: "b",
        minver: 1,
        webm: false,
        description: "Settings that might be needed by the decompressor. For Header Stripping (ContentCompAlgo=3), the bytes that were removed from the beggining of each frames of the track."
    },
    17026: {
        name: "DocType",
        level: 1,
        type: "s",
        mandatory: true,
        default: "matroska",
        minver: 1,
        description: "A string that describes the type of document that follows this EBML header. 'matroska' in our case or 'webm' for webm files."
    },
    17029: {
        name: "DocTypeReadVersion",
        level: 1,
        type: "u",
        mandatory: true,
        default: 1,
        minver: 1,
        description: "The minimum DocType version an interpreter has to support to read this file."
    },
    17030: {
        name: "EBMLVersion",
        level: 1,
        type: "u",
        mandatory: true,
        default: 1,
        minver: 1,
        description: "The version of EBML parser used to create the file."
    },
    17031: {
        name: "DocTypeVersion",
        level: 1,
        type: "u",
        mandatory: true,
        default: 1,
        minver: 1,
        description: "The version of DocType interpreter used to create the file."
    },
    // Many more EBML ID mappings follow...
};

// In TypeScript, we need to declare these modules
declare const require: any;

// Node.js module exports (converted from CommonJS to ES Module)
export const fs = require("fs");
export const stream = require("stream");
export const path = require("path");
export const events = require("events");
export const util = require("util");
export const url = require("url");
export const crypto = require("crypto");
export const buffer = require("buffer");
export const http = require("http");

// Export the inherits utility function
export const inherits = (() => {
    try {
        const utilModule = require("util");
        if (typeof utilModule.inherits !== "function") {
            throw new Error("util.inherits is not a function");
        }
        return utilModule.inherits;
    } catch (e) {
        return require("./utils/inherits");
    }
})(); 