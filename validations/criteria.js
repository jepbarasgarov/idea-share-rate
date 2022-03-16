var validation =  {
    "validator":{
        "$jsonSchema": {
            "bsonType": "object",
            "required": ["name", "create_ts", "update_ts"],
            "properties": {
                "name": {
                    "bsonType":    "string",
                    "maxLength":   256,
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


db.createCollection("criteria",validation)
