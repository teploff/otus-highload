<template>
  <div class="main">
    <section class="sign-in">
      <div class="container">
        <div class="md-layout md-gutter">
          <div class="md-layout-item sign-in-image">
            <figure>
              <img src="../../static/images/sign_in_image.jpg" alt="sing up image">
            </figure>

            <router-link class="sign-in-image-link" to="/sign-up">Create an account</router-link>
          </div>
          <div class="md-layout-item sign-in-form">
            <div class="md-layout-item register-form">
              <h2 class="form-title">Sign in</h2>

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
                </md-field>

                <div class="form-button">
                  <md-button
                    class="md-dense md-raised md-primary sign-in-button"
                    type="submit">
                    Log in
                  </md-button>
                </div>
              </form>
            </div>
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
import { required, email } from 'vuelidate/lib/validators';
import { apiUrl, headers } from '../const';

export default {
  name: 'SignIn',
  mixins: [validationMixin],
  data: () => ({
    payload: {
      email: String,
      password: String,
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
        this.signIn();
      }
    },
    signIn() {
      const path = `${apiUrl}/auth/sign-in`;
      axios.post(path, this.payload, { headers })
        .then((response) => {
          const tokenPair = JSON.parse(JSON.stringify(response.data));
          localStorage.setItem('access_token', tokenPair.access_token);
          localStorage.setItem('refresh_token', tokenPair.refresh_token);
          this.$router.push({ name: 'Questionnaires' });
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
.sign-in {
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

.sign-in-image {
  margin: 50px 10px 0 50px;
  width: 100%;
}

.sign-in-form {
  margin-left: 15px;
  margin-right: 15px;
}

.form-title {
  line-height: 1.66;
  font-weight: bold;
  color: #222;
  font-family: Poppins;
  font-size: 36px;
}

.register-form {
  margin-left: 15px;
  margin-right: 35px;
  padding-left: 34px;
}

.form-button {
  text-align: center;
}

.sign-in-button {
  width:40%;
}

.sign-in-image-link {
  font-size: 16px;
  display: block;
  text-align: center;
}
</style>
