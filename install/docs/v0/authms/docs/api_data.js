define({ "api": [
  {
    "type": "put",
    "url": "/:loginType/register",
    "title": "",
    "name": "Register",
    "group": "Auth",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "201": [
          {
            "group": "201",
            "optional": false,
            "field": "see",
            "description": "<p>See <a href=\"#api-Objects-User\">User</a>.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/handler.go",
    "groupTitle": "Auth"
  },
  {
    "type": "get",
    "url": "/status",
    "title": "",
    "name": "Status",
    "group": "Auth",
    "header": {
      "fields": {
        "Header": [
          {
            "group": "Header",
            "optional": false,
            "field": "x-api-key",
            "description": "<p>the api key</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "200": [
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Micro-service name.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "version",
            "description": "<p>Current running version.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "description",
            "description": "<p>Short description of the micro-service.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "canonicalName",
            "description": "<p>Canonical name of the micro-service.</p>"
          },
          {
            "group": "200",
            "type": "String",
            "optional": false,
            "field": "needRegSuper",
            "description": "<p>true if a super-user has been registered, false otherwise.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/handler.go",
    "groupTitle": "Auth"
  },
  {
    "type": "NULL",
    "url": "Device",
    "title": "",
    "name": "Device",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the device (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID for user who owns this device ID.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "deviceID",
            "description": "<p>The unique device ID string value.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date the device was created.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date the device was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/device.go",
    "groupTitle": "Objects"
  },
  {
    "type": "NULL",
    "url": "FacebookID",
    "title": "",
    "name": "FacebookID",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the facebook ID (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID for user who owns this facebook ID.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "facebookID",
            "description": "<p>The unique facebook ID string value.</p>"
          },
          {
            "group": "Success 200",
            "type": "Boolean",
            "optional": false,
            "field": "verified",
            "description": "<p>True if this login is verified, false otherwise.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date the facebook ID was inserted.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date the facebook ID value was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/facebook.go",
    "groupTitle": "Objects"
  },
  {
    "type": "NULL",
    "url": "Group",
    "title": "",
    "name": "Group",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the group (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>The unique group name string value.</p>"
          },
          {
            "group": "Success 200",
            "type": "Integer",
            "optional": false,
            "field": "accessLevel",
            "description": "<p>The access level for this group in (0 &gt;= accessLevel &lt;= 10)</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date the group was created.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date the group was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/group.go",
    "groupTitle": "Objects"
  },
  {
    "type": "NULL",
    "url": "User",
    "title": "",
    "name": "User",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the user (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "JWT",
            "description": "<p>JSON Web Token for accessing services.</p>"
          },
          {
            "group": "Success 200",
            "type": "Object",
            "optional": false,
            "field": "username",
            "description": "<p>See <a href=\"#api-Objects-Username\">Username</a>.</p>"
          },
          {
            "group": "Success 200",
            "type": "Object",
            "optional": false,
            "field": "phone",
            "description": "<p>See <a href=\"#api-Objects-VerifLogin\">VerifLogin</a>.</p>"
          },
          {
            "group": "Success 200",
            "type": "Object",
            "optional": false,
            "field": "email",
            "description": "<p>See <a href=\"#api-Objects-VerifLogin\">VerifLogin</a>.</p>"
          },
          {
            "group": "Success 200",
            "type": "Object",
            "optional": false,
            "field": "facebook",
            "description": "<p>See <a href=\"#api-Objects-FacebookID\">FacebookID</a>.</p>"
          },
          {
            "group": "Success 200",
            "type": "Object",
            "optional": false,
            "field": "group",
            "description": "<p>See <a href=\"#api-Objects-Group\">Group</a>.</p>"
          },
          {
            "group": "Success 200",
            "type": "Object",
            "optional": false,
            "field": "device",
            "description": "<p>See <a href=\"#api-Objects-Device\">Device</a>.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>Date the user was created.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>date the user was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/user.go",
    "groupTitle": "Objects"
  },
  {
    "type": "NULL",
    "url": "UserType",
    "title": "",
    "name": "UserType",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the userType (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "name",
            "description": "<p>Unique name of the user type.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date the user type was created.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date the user type was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/user_type.go",
    "groupTitle": "Objects"
  },
  {
    "type": "NULL",
    "url": "Username",
    "title": "",
    "name": "Username",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the username (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID for user who owns this Username.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "value",
            "description": "<p>The unique username string value.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date the username was created.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date the username was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/username.go",
    "groupTitle": "Objects"
  },
  {
    "type": "NULL",
    "url": "VerifLogin",
    "title": "",
    "name": "VerifLogin",
    "group": "Objects",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ID",
            "description": "<p>Unique ID of the verifiable login (can be cast to long Integer).</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "userID",
            "description": "<p>ID for user who owns this verifiable login.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "value",
            "description": "<p>The unique verifiable login string value.</p>"
          },
          {
            "group": "Success 200",
            "type": "Boolean",
            "optional": false,
            "field": "verified",
            "description": "<p>True if this login is verified, false otherwise.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created",
            "description": "<p>ISO8601 date the verifiable login was created.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lastUpdated",
            "description": "<p>ISO8601 date the verifiable login was last updated.</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "handler/http/verifiable_login.go",
    "groupTitle": "Objects"
  }
] });
