var validation =  {
    "validator":{
        "$jsonSchema": {
            "bsonType": "object",
            "required": ["name"],
            "properties": {
                "name": {
                    "bsonType":    "string",
                    "maxLength":   256,
                    "minLength":   1,
                    "description": "must be a string and is required",
                },
            },
        },
    },
        "validationLevel": "strict",
}


db.createCollection("position",validation)
