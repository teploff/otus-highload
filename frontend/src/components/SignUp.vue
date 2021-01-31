<template>
  <section class="form-elegant">
    <mdb-row>
      <mdb-col md="5">
        <mdb-card>
          <mdb-card-body class="mx-4">
            <form id="sign-up-form" novalidate @submit.prevent="checkForm">
              <div class="text-center">
                <h3 class="dark-grey-text mb-5"><strong>Sign up</strong></h3>
              </div>
              <div class="dark-grey-text text-left">
                <mdb-input v-model="signUpPayload.email" icon="envelope" label="Your email" type="email" required invalidFeedback="Please provide a valid email."/>
                <mdb-input v-model="signUpPayload.password" icon="lock" label="Your password" type="password" required invalidFeedback="Please provide a password."/>
                <mdb-input icon="unlock-alt" label="Type your password again" type="password" required invalidFeedback="Please provide a repeated password."/>
                <mdb-input v-model="signUpPayload.name" icon="user-ninja" label="Your name" type="text" required invalidFeedback="Please provide a valid name."/>
                <mdb-input v-model="signUpPayload.surname" icon="user-ninja" label="Your surname" type="text" required invalidFeedback="Please provide a valid surname."/>
                <mdb-date-picker v-model="signUpPayload.birthday" icon="birthday-cake" label="Your birthday" disabledFuture disableClear invalidFeedback="Please provide a valid birthday."/>
                <mdb-form-inline>
                  <mdb-input type="radio" id="option5-1" name="groupOfMaterialRadios" radioValue="male" v-model="signUpPayload.sex" label="Male"/>
                  <mdb-input type="radio" id="option5-2" name="groupOfMaterialRadios" radioValue="female" v-model="signUpPayload.sex" label="Female" />
                </mdb-form-inline>
                <mdb-input v-model="signUpPayload.city" icon="city" label="Your city" type="text" required invalidFeedback="Please provide a city."/>
                <mdb-input v-model="signUpPayload.interests" icon="camera-retro" type="textarea" label="Your interests" required invalidFeedback="Please provide some interests."/>
              </div>
              <div class="text-center mb-3">
                <mdb-btn type="submit" gradient="blue" rounded class="btn-block z-depth-1a">Sign up</mdb-btn>
              </div>
            </form>
          </mdb-card-body>
          <mdb-modal-footer class="mx-5 pt-3 mb-1">
            <p class="font-small grey-text d-flex justify-content-end"> Have an account?<router-link to="/sign-in"> Sign In</router-link></p>
          </mdb-modal-footer>
        </mdb-card>
      </mdb-col>
    </mdb-row>
  </section>
</template>

<script>
import {mdbRow, mdbCol, mdbCard, mdbCardBody, mdbInput, mdbBtn, mdbModalFooter, mdbDatePicker, mdbFormInline} from "mdbvue";
import {signUpUser} from "@/api/auth.api";
const moment = require('moment');

export default {
  components: {
    mdbRow,
    mdbCol,
    mdbCard,
    mdbCardBody,
    mdbBtn,
    mdbModalFooter,
    mdbInput,
    mdbDatePicker,
    mdbFormInline
  },
  data: () => ({
    showModal: false,
    signUpPayload: {
      email: '',
      password: '',
      repeatedPassword: '',
      name: '',
      surname: '',
      birthday: moment().subtract(18, 'years').format('YYYY-MM-DD'),
      sex: 'male',
      city: '',
      interests: '',
    }
  }),
  name: "SignUp",
  methods: {
    checkForm(event) {
      event.target.classList.add('was-validated');

      const form = document.getElementById('sign-up-form');
      if (form.checkValidity()) {
        this.signUp();
      }
    },
    async signUp() {
      try {
        let request = this.signUpPayload
        request.birthday = moment(this.signUpPayload.birthday, 'YYYY-MM-DD').utc().toISOString()

        await signUpUser(request)
        this.$notify.success({message: 'Welcome to Social Network!\nWe are glad you are with us! You can start using our ' +
              'platform.', position: 'top right', timeOut: 10000});
        await this.$router.push({name: 'SignIn'})
      } catch (error) {
          this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
      }
    },
  }
}
</script>

<style scoped>

</style>