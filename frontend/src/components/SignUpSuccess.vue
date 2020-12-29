<template>
  <div class="main">
    <section class="sign-up-success">
      <div class="container">
        <div class="md-layout md-gutter">
          <div class="md-layout-item">
            <h2 class="success-title">You have signed up successfully</h2>
            <hr class="separate-line">
            <div class="success-text">
              Welcome to Social Network!
              We are glad you are with us! You can start using our platform.
              Please, proceed to a login page.
            </div>
            <div class="form-button">
              <md-button
                class="md-dense md-raised md-primary sign-up-success-button"
                v-on:click="signIn">
                Sign In
              </md-button>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: 'SignUpSuccess',
  methods: {
    signIn() {
      this.$router.push({ name: 'SignIn' });
    },
  },
};
</script>

<style scoped>
.main {
  background: #f8f8f8;
  padding: 250px 0;
}
.sign-up-success {
  margin-top: 75px;
  margin-bottom: 75px;
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

.success-title {
  line-height: 1.66;
  font-weight: bold;
  color: #222;
  font-family: Poppins;
  font-size: 36px;
  text-align: center;
}

.separate-line {
  margin-left: 40px;
  margin-right: 40px;
}

.success-text {
  margin-top: 10px;
  margin-left: 50px;
  margin-right: 50px;
  font-family: Poppins;
  font-size: 20px;
  text-align: center;
}

.form-button {
  text-align: center;
  margin-top: 30px;
  margin-bottom: 30px;
}

.sign-up-success-button {
  width:40%;
}
</style>
