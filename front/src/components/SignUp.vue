<template>
  <div class="main">
    <section class="sign-up">
      <div class="container">
        <div class="sign-up-content">
          <div class="md-layout-item md-size-337">
            <h2 class="form-title">Sign up</h2>
            <div class="sign-up-form">
              <form class="register-form">
                <md-field class="form-group">
                  <md-icon>email</md-icon>
                  <label>Your Email</label>
                  <md-input v-model="payload.email" type="email" required></md-input>
                </md-field>

                <md-field class="form-group">
                  <md-icon>lock</md-icon>
                  <label>Your Password</label>
                  <md-input v-model="payload.password" type="password" required></md-input>
                </md-field>

                <md-field class="form-group" :md-toggle-password="false">
                  <md-icon>lock_open</md-icon>
                  <label>Type your password again</label>
                  <md-input type="password" required></md-input>
                </md-field>

                <md-field class="form-group">
                  <md-icon>person</md-icon>
                  <label>Your Name</label>
                  <md-input v-model="payload.name" type="text" required></md-input>
                </md-field>

                <md-field class="form-group">
                  <md-icon>person</md-icon>
                  <label>Your Surname</label>
                  <md-input v-model="payload.surname" type="text" required></md-input>
                </md-field>

                <div class="form-group">
                  <md-datepicker v-model="payload.birthday" required>
                    <label>Your Birthday</label>
                  </md-datepicker>
                </div>

                <div class="form-group">
                  <md-radio v-model="payload.sex" value="male">Male</md-radio>
                  <md-radio v-model="payload.sex" value="female">Female</md-radio>
                </div>

                <md-field class="form-group">
                  <md-icon>location_city</md-icon>
                  <label>Your City</label>
                  <md-input v-model="payload.city" type="text" required></md-input>
                </md-field>

                <div class="form-group">
                  <md-icon>flight_takeoff</md-icon>
                  <md-field>
                    <label>Your Interests</label>
                    <md-textarea v-model="payload.interests" required></md-textarea>
                  </md-field>
                </div>
                <div class="form-group form-button">
                  <md-button class="md-raised" id="signUpButton" v-on:click="signUp">
                    Register
                  </md-button>
                </div>
              </form>
            </div>
          </div>
          <div class="md-layout-item md-size-337">
            <div class="sign-up-image">
              <table width="100%">
                <tr>
                  <td id="signup_image_row_1"></td>
                </tr>
                <tr>
                  <td>
                    <figure>
                      <img src="../../static/images/sign_up_image.jpg" alt="sing up image">
                    </figure>
                  </td>
                </tr>
                <tr>
                  <td>
                    <a href="#" class="signup-image-link">I am already member</a>
                  </td>
                </tr>
                <tr>
                  <td id='signup_image_row_4'></td>
                </tr>
              </table>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
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
