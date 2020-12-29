<template>
  <div>
    <div class="logout-button">
      <md-button
        class="md-dense md-raised md-primary"
        @click="logOut">
        Log out
      </md-button>
    </div>
    <div v-show="cards.count === 0">
      <md-empty-state
        v-show="cards.count === 0"
        md-rounded
        md-icon="star"
        md-label="You are the first!"
        md-description="Congratulation you are the first user in this portal!
        While there aren't other people but they will appear soon"
        style="width: 600px; height: 600px; position: center">
      </md-empty-state>
    </div>
    <md-table v-show="cards.count !== 0">
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
            </md-card-header>

            <md-card-expand>
              <md-card-actions md-alignment="right">
                <md-card-expand-trigger>
                  <md-button class="learn-more-button" style="color: #337ab7">Learn more</md-button>
                </md-card-expand-trigger>
              </md-card-actions>

              <md-card-expand-content>
                <md-card-content>
                  <p style="text-align: left"> <b>Sex:</b> {{ card.sex }} </p>
                  <p style="text-align: left">
                    <b>Birthday:</b> {{ new Date(card.birthday) | dateFormat('DD.MM.YYYY') }}
                  </p>
                  <p style="text-align: left"> <b>City:</b> {{ card.city }} </p>
                  <p style="text-align: left"> <b>Interests:</b> {{ card.interests }} </p>
                </md-card-content>
              </md-card-expand-content>
            </md-card-expand>
          </md-card>
        </div>
      </md-table-row>
      <md-table-row class="pagination-row">
        <div>
          <paginate
            v-model="page"
            :page-count="Math.ceil(cards.count / countCardsInWindow)"
            :click-handler="paginatorClick"
            :prev-text="'Prev'"
            :next-text="'Next'"
            :container-class="'pagination'"
            :page-class="'page-item'"
            :first-last-button="true"
          >
          </paginate>
        </div>
      </md-table-row>
    </md-table>
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
    countCardsInWindow: 10,
    page: 1,
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
    paginatorClick(pageNum) {
      const path = `${apiUrl}/questionnaires`;
      headers.Authorization = localStorage.getItem('access_token');

      this.payload.offset = (pageNum - 1) * this.countCardsInWindow;

      axios.post(path, this.payload, { headers })
        .then((response) => {
          this.cards = JSON.parse(JSON.stringify(response.data));
        })
        .catch((error) => {
          const err = JSON.parse(JSON.stringify(error.response));
          if (err.status === 401) {
            this.refreshToken();
          }
        });
    },
    logOut() {
      localStorage.clear();
      this.$router.push({ name: 'SignIn' });
    },
  },
  created() {
    this.getQuestionnaires();
  },
};
</script>

<style scoped>
.card-expansion {
  margin: 0 175px 10px 175px;
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

.md-empty-state {
  max-width: 600px;
}

.logout-button {
  padding: 15px 25px 0 0;
  text-align: right;
}

.learn-more-button{
  font-weight: bolder;
}

</style>
