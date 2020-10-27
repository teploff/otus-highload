# User generation
<img align="right" width="160" src="./static/image.png">

##  What's the point of that?
This tool is based on [Faker](https://github.com/joke2k/faker) lib. It cans without troubles to generate users which consist of fields, such as:
- email;
- password;
- name;
- surname;
- birthday;
- sex;
- city;
- interests.

After generation users are stored in *.txt file for future actions.

## Using
### Requirements
- Python >= 3.6 version.
### Launching
```shell script
python main.py --count=<count of generated users> -path=<path where users should stored> 
```

Example:
```shell script
python main.py --count=100 -path=.
```
It means that there will be generate 100 users which will be sored in directory "."