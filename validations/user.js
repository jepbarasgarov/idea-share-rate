var validation =  {
    "validator":{
        "$jsonSchema": {
            "bsonType": "object",
            "required": ["firstname", "lastname", "username", "create_ts", "update_ts","role", "status", "password"],
            "properties": {
                "username": {
                    "bsonType":    "string",
                    "maxLength":   32,
                    "minLength":   1,
                    "description": "must be a string and is required",
                },
                "firstname": {
                    "bsonType":    "string",
                    "maxLength":   32,
                    "minLength":   1,
                    "description": "must be a string and is required",
                },
                "lastname":{
                    "bsonType":    "string",
                    "maxLength":   32,
                    "minLength":   1,
                    "description": "must be a string and is required",
                },
                "role":{
                    "enum": ["ADMIN", "USER"],
                    "description": "must be a ADMIN or USER and is required",
                },
                "status":{
                    "enum": ["ACTIVE", "BLOCKED"],
                    "description": "must be a ACTIVE or BLOCKED and is required",
                },
                "password":{
                    "bsonType": "string",
                    "description": "must be a ACTIVE or BLOCKED and is required",
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


db.createCollection("user",validation)



db.user.insertOne({
    "username":"admin",
    "firstname":"admin",
    "lastname":"admin",
    "role":"ADMIN",
    "status":"ACTIVE",
    "password":"$2a$12$BWKzPWJXBqhHuvhKlCCYQ.tbmBCJq3IULfqBPjMHSrZeJnLxk/Fhq",
    "create_ts": new Date(),
    "update_ts":  new Date()
})





