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


db.createCollection("genre",validation)
db.genre.createIndex({"name":1}, {unique:true})

