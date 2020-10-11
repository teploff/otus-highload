<template>
  <form>
    <md-field>
      <md-icon>email</md-icon>
      <label>Your Email</label>
      <md-input v-model="payload.email" type="email" required></md-input>
    </md-field>

    <md-field>
      <md-icon>lock</md-icon>
      <label>Your Password</label>
      <md-input v-model="payload.password" type="password" required></md-input>
    </md-field>

    <md-field :md-toggle-password="false">
      <md-icon>lock_open</md-icon>
      <label>Type your password again</label>
      <md-input type="password" required></md-input>
    </md-field>

    <md-field>
      <md-icon>person</md-icon>
      <label>Your Name</label>
      <md-input v-model="payload.name" type="text" required></md-input>
    </md-field>

    <md-field>
      <md-icon>person</md-icon>
      <label>Your Surname</label>
      <md-input v-model="payload.surname" type="text" required></md-input>
    </md-field>

    <div>
      <md-datepicker v-model="payload.birthday" required>
        <label>Your Birthday</label>
      </md-datepicker>
    </div>

    <div>
      <md-radio v-model="payload.sex" value="male">Male</md-radio>
      <md-radio v-model="payload.sex" value="female">Female</md-radio>
    </div>

    <md-field>
      <md-icon>location_city</md-icon>
      <label>Your City</label>
      <md-input v-model="payload.city" type="text" required></md-input>
    </md-field>

    <md-icon>flight_takeoff</md-icon>
    <md-field>
      <label>Your Interests</label>
      <md-textarea v-model="payload.interests" required></md-textarea>
    </md-field>

    <md-button class="md-raised" id="signUpButton" v-on:click="signUp">Register</md-button>
  </form>
</template>

<script>
import axios from 'axios';

export default {
  name: 'SignUp',
  data: () => ({
    payload: {
      email: '',
      password: '',
      name: '',
      surname: '',
      birthday: Date(),
      sex: 'male',
      city: '',
      interests: '',
    },
    error: null,
    info: null,
  }),
  methods: {
    signUp() {
      const path = 'http://localhost:9999/auth/sign-up';
      const headers = {
        'Content-Type': 'application/json',
      };
      axios.post(path, this.payload, { headers })
        .then((response) => {
          this.info = response.data;
        })
        .catch((error) => {
          // eslint-отключение следующей строки
          this.error = JSON.parse(JSON.stringify(error.response));
          console.log(error);
        });
      console.log(this.error);
    },
  },
};
</script>
