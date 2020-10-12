<template>
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
}

.md-card {
  width: 250px;
  margin: 4px;
  display: inline-block;
  vertical-align: top;
}
</style>
