<template>
  <section class="form-elegant">
    <mdb-row>
      <mdb-col md="5">
        <mdb-card>
          <mdb-card-body class="mx-4">
            <form id="sign-in-form" novalidate @submit.prevent="checkForm">
              <div class="text-center">
                <h3 class="dark-grey-text mb-5"><strong>Sign in</strong></h3>
              </div>
              <div class="dark-grey-text text-left">
                <mdb-input v-model="email" icon="envelope" label="Your email" type="email" required invalidFeedback="Please provide a valid email."/>
                <mdb-input v-model="password" icon="lock" label="Your password" type="password" required invalidFeedback="Please provide a password."/>
              </div>
              <div class="text-center mb-3">
                <mdb-btn type="submit" gradient="blue" rounded class="btn-block z-depth-1a">Sign in</mdb-btn>
              </div>
            </form>
          </mdb-card-body>
          <mdb-modal-footer class="mx-5 pt-3 mb-1">
            <p class="font-small grey-text d-flex justify-content-end">Not a member?<router-link to="/sign-up">Sign Up</router-link></p>
          </mdb-modal-footer>
        </mdb-card>
      </mdb-col>
    </mdb-row>
  </section>
</template>

<script>
import { mdbRow, mdbCol, mdbCard, mdbCardBody, mdbInput, mdbBtn, mdbModalFooter } from 'mdbvue';
import {signInUser} from "@/api/auth.api";

export default {
  name: "SignIn",
  components: {
    mdbRow,
    mdbCol,
    mdbCard,
    mdbCardBody,
    mdbInput,
    mdbBtn,
    mdbModalFooter
  },
  data: () => ({
    showModal: false,
    email: '',
    password: '',
  }),
  methods: {
    checkForm(event) {
      event.target.classList.add('was-validated');

      const form = document.getElementById('sign-in-form');
      if (form.checkValidity()) {
        this.signIn()
      }
    },
    async signIn() {
      try {
        const response = await signInUser(this.email, this.password)

        localStorage.setItem('accessToken', response.data.accessToken)
        localStorage.setItem('refreshToken', response.data.refreshToken)

        await this.$router.push({name: 'Home'})
      } catch (error) {
        this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
      }
    },
  }
}
</script>

<style scoped>

</style>