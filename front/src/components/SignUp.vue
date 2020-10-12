<template>
  <div class="main">
    <section class="sign-up">
      <div class="container">
        <div class="md-layout md-gutter">
          <div class="md-layout-item register-form">
            <h2 class="form-title">Sign up</h2>

            <form>
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

          <div class="md-layout-item sign-up-image">
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

<style>
  .main {
    background: #f8f8f8;
    padding: 150px 0;
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
    margin: 0 55px;
  }

  .form-title {
    line-height: 1.66;
    font-weight: bold;
    color: #222;
    font-family: Poppins;
    font-size: 36px;
  }

  .register-form {
    margin-left: 75px;
    margin-right: 75px;
    padding-left: 34px;
  }

  .sign-up-form, .sign-up-image, .sign-in-form, .sign-in-image {
    width: 50%;
    overflow: hidden;
  }

  .form-group:last-child {
    margin-bottom: 0;
  }


  .form-button {
    text-align: center;
  }
</style>
