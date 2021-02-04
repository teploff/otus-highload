```shell script
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890", "name": "Bob", "surname": "Tallor", "birthday": "1994-04-10T20:21:25+00:00", "sex": "male", "city": "New Yourk", "interests": "programming"}' \
    http://localhost:10000/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "alice@email.com", "password": "1234567890", "name": "Alice", "surname": "Swift", "birthday": "1995-10-10T20:21:25+00:00", "sex": "female", "city": "California", "interests": "running"}' \
    http://localhost:10000/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "henry@email.com", "password": "1234567890", "name": "Henry", "surname": "Cavill", "birthday": "1993-08-19T20:21:25+00:00", "sex": "male", "city": "Washington", "interests": "sport"}' \
    http://localhost:10000/auth/sign-up
```

```shell script
docker exec -it auth-storage mysql -uroot -ppassword auth
```

```shell script
update user set id = 'f14d517b-e1a6-4dc1-940a-0657185b4391' where email = 'bob@email.com';
update user set id = 'e1f46383-0a1e-49db-bf79-f0eceab6427c' where email = 'alice@email.com';
update user set id = '461f8a43-252a-4b1b-8077-96df252e73cb' where email = 'henry@email.com';
```


```shell script
export BOB_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
export ALICE_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "alice@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
export HENRY_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "henry@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
```

```shell script
export ALICE_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    http://localhost:9999/auth/user?email=alice@email.com | jq -r '.user_id')
```
    
```shell script
export ALICE_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    http://localhost:9999/auth/user?email=alice@email.com | jq -r '.user_id')
```

```shell script
export ALICE_ID=e1f46383-0a1e-49db-bf79-f0eceab6427c
```

```shell script
export CHAT_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    -d '{"companion_id": "'"$ALICE_ID"'"}' \
    http://localhost:10000/messenger/create-chat | jq -r '.chat_id')
```

```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${ALICE_ACCESS_TOKEN}" \
http://localhost:10000/messenger/messages?chat_id=$CHAT_ID
```

    
