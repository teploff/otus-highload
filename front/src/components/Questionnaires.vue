<template>
  <md-table>
    <md-table-row>
      <div class="card-expansion">
        <md-card v-for="card in cards.questionnaires" v-bind:key="card.email">
          <md-card-media>
            <img v-if="card.sex === 'male'" src="static/images/boy.png" alt="People">
            <img v-else src="static/images/girl.png" alt="People">
          </md-card-media>

          <md-card-header>
            <div class="md-title">{{ card.name }} {{ card.surname }}</div>
            <div class="md-subhead">{{ card.email }}</div>
            <div class="md-subhead">{{ card.sex }}</div>
          </md-card-header>

          <md-card-expand>
            <md-card-actions md-alignment="right">
              <md-card-expand-trigger>
                <md-button>Learn more</md-button>
              </md-card-expand-trigger>
            </md-card-actions>

            <md-card-expand-content>
              <md-card-content>
                <p> Birthday: {{ card.birthday }} </p>
                <p> City: {{ card.city }} </p>
                <p> Interests: {{ card.interests }} </p>
              </md-card-content>
            </md-card-expand-content>
          </md-card-expand>
        </md-card>
      </div>
    </md-table-row>
    <md-table-row class="pagination-row">
      todo need paginator
    </md-table-row>
  </md-table>
</template>

<script>
import axios from 'axios';
import { apiUrl, headers } from '../const';

export default {
  name: 'Questionnaires',
  data: () => ({
    payload: {
      limit: 10,
      offset: 0,
    },
    cards: {
      questionnaires: null,
      count: 0,
    },
  }),
  methods: {
    getQuestionnaires() {
      const path = `${apiUrl}/questionnaires`;
      headers.Authorization = localStorage.getItem('access_token');

      axios.post(path, this.payload, { headers })
        .then((response) => {
          this.cards = JSON.parse(JSON.stringify(response.data));
        })
        .catch((error) => {
          const err = JSON.parse(JSON.stringify(error.response));
          if (err.status === 401) {
            this.refreshToken();
          }
          console.log(err);
        });
    },
    refreshToken() {
      const path = `${apiUrl}/auth/token`;
      const refreshToken = localStorage.getItem('refresh_token');

      if (refreshToken === null) {
        this.$router.push({ name: 'SignIn' });
      }

      const payload = {
        refresh_token: refreshToken,
      };
      axios.put(path, payload)
        .then((response) => {
          const tokenPair = JSON.parse(JSON.stringify(response.data));
          localStorage.setItem('access_token', tokenPair.access_token);
          localStorage.setItem('refresh_token', tokenPair.refresh_token);
          this.$router.push({ name: 'Questionnaires' });
        })
        .catch((error) => {
          const err = JSON.parse(JSON.stringify(error.response));
          if (err.status === 401) {
            this.$router.push({ name: 'SignIn' });
          }
          console.log(err);
        });
    },
  },
  created() {
    this.getQuestionnaires();
  },
};
</script>

<style scoped>
.card-expansion {
  margin: 25px 150px 10px 150px;
  text-align: center;
}

.md-card {
  width: 250px;
  margin: 4px;
  display: inline-block;
  vertical-align: top;
}

.pagination-row {
  text-align: center;
  margin-bottom: 75px;
}
</style>
