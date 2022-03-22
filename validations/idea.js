var validationIdea =  {
    "validator":{
        "$jsonSchema": {
            "bsonType": "object",
            "required": ["name", "worker", "date", "genre", "mechanics", "links", "description", "files", "rates", "create_ts", "update_ts"],
            "properties": {
                "name": {
                    "bsonType":    "string",
                    "maxLength":   256,
                    "minLength":   1,
                    "description": "must be a string and is required",
                },
                "genre":{
                    "bsonType":    "string",
                    "maxLength":   256,
                    "minLength":   1,
                    "description": "must be a string and is required",
                },
                "worker": {
                    "bsonType":    "object",
                    "required": ["_id", "firstname", "lastname", "position"],
                    "properties": {
                        "_id":{
                            "bsonType":   "objectId",
                        },
                        "firstname": {
                            "bsonType":    "string",
                            "maxLength":   64,
                            "minLength":   1,
                            "description": "must be a string and is required",
                        },
                        "lastname":{
                            "bsonType":    "string",
                            "maxLength":   64,
                            "minLength":   1,
                            "description": "must be a string and is required",
                        },
                        "position": {
                            "bsonType":    "string",
                            "maxLength":   128,
                            "minLength":   1,
                            "description": "must be a string and is required",
                        },
                    },
                },
                "date":{
                    "bsonType":    "date",
                    "description": "must be a date and is required",

                },
                "mechanics":{
                    "bsonType":    "array",
                    "minItems":1,
                    "description": "must be a array of strings and is required",
                    "items":{
                        "bsonType": "string",
                    }

                },
                "links":{
                    "bsonType":    "array",
                    "items":{
                        "bsonType": "object",
                        "required": ["label", "url"],
                        "properties":{
                            "label":{"bsonType":"string"},
                            "url":{"bsonType":"string"}
                        }
                    }
                    
                },
                "description":{
                    "bsonType":"string"
                },
                "files":{
                    "bsonType":    "array",
                    "items":{
                        "bsonType": "object",
                        "required": ["sketch_id", "name", "file_path"],
                        "properties":{
                            "sketch_id":{"bsonType":"objectId"},
                            "name":{"bsonType":"string"},
                            "file_path":{"bsonType":"string"}
                        }
                    }
                },
                "rates":{
                    "bsonType":    "array",
                    "items":{
                        "bsonType": "object",
                        "required": ["criteria_id", "user_id", "rate"],
                        "properties":{
                            "criteria_id":{"bsonType":"objectId"},
                            "user_id":{"bsonType":"objectId"},
                            "rate":{
                                "bsonType":"int",
                                "maximum":10,
                                "minimum":0
                            }
                        }
                    }
                },
                "create_ts":{
                    "bsonType":    "date",
                },
                "update_ts":{
                    "bsonType":    "date",
                }

            },
        },
    },
        "validationLevel": "strict",
}


db.createCollection("idea",validationIdea)
db.idea.createIndex({"worker._id":1})
db.idea.createIndex({"date":1})



