<template>
  <div class="main">
    <section class="sign-up">
      <div class="container">
        <div class="md-layout md-gutter">
          <div class="md-layout-item register-form">
            <h2 class="form-title">Sign up</h2>

            <form novalidate @submit.prevent="makeValidate">
              <md-field :class="getValidationClass('email')">
                <md-icon>email</md-icon>
                <label>Your Email</label>
                <md-input v-model="payload.email" type="email"></md-input>
                <span class="md-error" v-if="!$v.payload.email.required">Email is required</span>
                <span class="md-error" v-else-if="!$v.payload.email.email">Invalid email</span>
              </md-field>

              <md-field :class="getValidationClass('password')">
                <md-icon>lock</md-icon>
                <label>Your Password</label>
                <md-input v-model="payload.password" type="password"></md-input>
                <span class="md-error" v-if="!$v.payload.password.required">
                  Password is required
                </span>
                <span class="md-error" v-else-if="!$v.payload.password.minLength">
                  Password should be more than 6 characters
                </span>
                <span class="md-error" v-else-if="!$v.payload.password.maxLength">
                  Password should be less than 20 characters
                </span>
              </md-field>

              <md-field
                :md-toggle-password="false"
                :class="getValidationClass('repeatedPassword')">
                <md-icon>lock_open</md-icon>
                <label>Type your password again</label>
                <md-input v-model="payload.repeatedPassword" type="password"></md-input>
                <span class="md-error" v-if="!$v.payload.repeatedPassword.required">
                  Repeated password is required
                </span>
                <span class="md-error" v-else-if="!$v.payload.repeatedPassword.sameAs">
                  Passwords should are same
                </span>
              </md-field>

              <md-field :class="getValidationClass('name')">
                <md-icon>person</md-icon>
                <label>Your Name</label>
                <md-input v-model="payload.name" type="text"></md-input>
                <span class="md-error" v-if="!$v.payload.name.required">
                  Name is required
                </span>
                <span class="md-error" v-else-if="!$v.payload.name.minLength">
                  Name should be more than 2 characters
                </span>
                <span class="md-error" v-else-if="!$v.payload.name.maxLength">
                  Name should be less than 20 characters
                </span>
              </md-field>

              <md-field :class="getValidationClass('surname')">
                <md-icon>person</md-icon>
                <label>Your Surname</label>
                <md-input v-model="payload.surname" type="text"></md-input>
                <span class="md-error" v-if="!$v.payload.surname.required">
                  Surname is required
                </span>
                <span class="md-error" v-else-if="!$v.payload.surname.minLength">
                  Surname should be more than 6 characters
                </span>
                <span class="md-error" v-else-if="!$v.payload.surname.maxLength">
                  Surname should be less than 20 characters
                </span>
              </md-field>

              <div>
                <md-datepicker v-model="payload.birthday">
                  <label>Your Birthday</label>
                </md-datepicker>
              </div>

              <div class="gender-item">
                <md-radio v-model="payload.sex" value="male">Male</md-radio>
                <md-radio v-model="payload.sex" value="female">Female</md-radio>
              </div>

              <md-field :class="getValidationClass('city')">
                <md-icon>location_city</md-icon>
                <label>Your City</label>
                <md-input v-model="payload.city" type="text"></md-input>
                <span class="md-error" v-if="!$v.payload.city.required">
                  City is required
                </span>
                <span class="md-error" v-else-if="!$v.payload.city.minLength">
                  City should be more than 2 characters
                </span>
                <span class="md-error" v-else-if="!$v.payload.city.maxLength">
                  City should be less than 20 characters
                </span>
              </md-field>

              <md-field :class="getValidationClass('interests')">
                <md-icon>flight_takeoff</md-icon>
                <label>Your Interests</label>
                <md-textarea v-model="payload.interests"></md-textarea>
                <span class="md-error" v-if="!$v.payload.interests.required">
                Interests is required
              </span>
                <span class="md-error" v-else-if="!$v.payload.interests.minLength">
                Interests should be more than 2 characters
              </span>
              </md-field>

              <div class="form-button">
                <md-button
                  class="md-dense md-raised md-primary sign-up-button"
                  id="signUpButton"
                  type="submit">
                  Register
                </md-button>
              </div>
            </form>
          </div>

          <div class="md-layout-item sign-up-image">
            <figure>
              <img src="../assets/sign_up_image.jpg" alt="sing up image">
            </figure>

            <router-link class="signup-image-link" to="/sign-in">I am already member</router-link>
          </div>
        </div>
      </div>
    </section>
    <FlashMessage :position="'right top'"></FlashMessage>
  </div>
</template>

<script>
import axios from 'axios';
import { validationMixin } from 'vuelidate';
import { required, email, minLength, maxLength, sameAs } from 'vuelidate/lib/validators';
import { apiUrl, headers } from '@/const';

export default {
  name: 'SignUp',
  mixins: [validationMixin],
  data: () => ({
    payload: {
      email: null,
      password: '',
      repeatedPassword: '',
      name: '',
      surname: '',
      birthday: new Date(),
      sex: 'male',
      city: '',
      interests: '',
    },
  }),
  validations: {
    payload: {
      email: {
        required,
        email,
      },
      password: {
        required,
        minLength: minLength(6),
        maxLength: maxLength(20),
      },
      repeatedPassword: {
        required,
        sameAs: sameAs('password'),
      },
      name: {
        required,
        minLength: minLength(2),
        maxLength: maxLength(20),
      },
      surname: {
        required,
        minLength: minLength(2),
        maxLength: maxLength(20),
      },
      city: {
        required,
        minLength: minLength(2),
        maxLength: maxLength(20),
      },
      interests: {
        required,
        minLength: minLength(2),
      },
    },
  },
  methods: {
    getValidationClass(fieldName) {
      const field = this.$v.payload[fieldName];

      if (field) {
        return {
          'md-invalid': field.$invalid && field.$dirty,
        };
      }

      return {
        'md-invalid': false,
      };
    },
    makeValidate() {
      this.$v.$touch();

      if (!this.$v.$invalid) {
        this.signUp();
      }
    },
    signUp() {
      const path = `${apiUrl}/auth/sign-up`;
      axios.post(path, this.payload, { headers })
        .then(() => {
          this.$router.push({ name: 'SignUpSuccess' });
        })
        .catch((error) => {
          const err = JSON.parse(JSON.stringify(error.response));
          this.flashMessage.error(
            { title: 'Error Message Title',
              message: err.data.message,
              position: 'center',
              icon: '../../static/images/error.svg',
            });
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
