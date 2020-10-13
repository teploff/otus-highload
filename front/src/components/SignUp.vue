<template>
  <div class="main">
    <section class="sign-up">
      <div class="container">
        <div class="md-layout md-gutter">
          <div class="md-layout-item register-form">
            <h2 class="form-title">Sign up</h2>

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

              <div class="gender-item">
                <md-radio v-model="payload.sex" value="male">Male</md-radio>
                <md-radio v-model="payload.sex" value="female">Female</md-radio>
              </div>

              <md-field>
                <md-icon>location_city</md-icon>
                <label>Your City</label>
                <md-input v-model="payload.city" type="text" required></md-input>
              </md-field>

              <div>
                <md-icon>flight_takeoff</md-icon>
                <md-field>
                  <label>Your Interests</label>
                  <md-textarea v-model="payload.interests" required></md-textarea>
                </md-field>
              </div>

              <div class="form-button">
                <md-button
                  class="md-dense md-raised md-primary sign-up-button"
                  id="signUpButton"
                  v-on:click="signUp">
                  Register
                </md-button>
              </div>
            </form>
          </div>

          <div class="md-layout-item sign-up-image">
            <figure>
              <img src="../../static/images/sign_up_image.jpg" alt="sing up image">
            </figure>

            <router-link class="signup-image-link" to="/sign-in">I am already member</router-link>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios';
import { apiUrl, headers } from '../const';

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
  }),
  methods: {
    signUp() {
      const path = `${apiUrl}/auth/sign-up`;
      axios.post(path, this.payload, { headers })
        .then(() => {
          this.$router.push({ name: 'SignUpSuccess' });
        })
        .catch((error) => {
          const err = JSON.parse(JSON.stringify(error.response));
          console.log(err);
        });
    },
  },
};
</script>

<style scoped>
  .main {
    background: #f8f8f8;
    padding: 150px 0;
    font-family: Poppins,serif;
  }
  .sign-up {
    margin-bottom: 150px;
  }

  .container {
    width: 900px;
    background: #fff;
    margin: 0 auto;
    box-shadow: 0 15px 17px 0 rgba(0, 0, 0, 0.05);
    -moz-box-shadow: 0 15px 17px 0 rgba(0, 0, 0, 0.05);
    -webkit-box-shadow: 0 15px 17px 0 rgba(0, 0, 0, 0.05);
    -o-box-shadow: 0 15px 17px 0 rgba(0, 0, 0, 0.05);
    -ms-box-shadow: 0px 15px 17px 0px rgba(0, 0, 0, 0.05);
    border-radius: 20px;
    -moz-border-radius: 20px;
    -webkit-border-radius: 20px;
    -o-border-radius: 20px;
    -ms-border-radius: 20px;
  }

  .md-layout {
    padding: 75px 0;
  }

  .sign-up-image {
    margin: 250px 10px 0 10px;
    width: 100%;
  }

  .form-title {
    line-height: 1.66;
    font-weight: bold;
    color: #222;
    font-size: 36px;
  }

  .register-form {
    margin-left: 75px;
    margin-right: 75px;
    padding-left: 34px;
  }

  .gender-item {
    text-align: center;
  }

  .signup-image-link {
    font-size: 16px;
    display: block;
    text-align: center;
  }

  .form-button {
    text-align: center;
  }

  .sign-up-button {
    width:40%;
  }
</style>
