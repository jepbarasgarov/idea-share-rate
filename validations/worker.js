
var validation =  {
    "validator":{
        "$jsonSchema": {
            "bsonType": "object",
            "required": ["firstname", "lastname", "position", "create_ts", "update_ts"],
            "properties": {
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


db.createCollection("worker",validation)

