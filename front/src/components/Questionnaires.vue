<template>
  <div class="card-expansion">
    <md-card>
      <md-card-media>
        <img src="static/images/boy.png" alt="People">
      </md-card-media>

      <md-card-header>
        <div class="md-title">Title goes here</div>
        <div class="md-subhead">Subtitle here</div>
      </md-card-header>

      <md-card-expand>
        <md-card-actions md-alignment="space-between">
          <div>
            <md-button>Action</md-button>
            <md-button>Action</md-button>
          </div>

          <md-card-expand-trigger>
            <md-button class="md-icon-button">
              <md-icon>keyboard_arrow_down</md-icon>
            </md-button>
          </md-card-expand-trigger>
        </md-card-actions>

        <md-card-expand-content>
          <md-card-content>
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Optio itaque ea, nostrum odio.
            Dolores, sed accusantium quasi non, voluptas eius illo quas,
            saepe voluptate pariatur in deleniti minus sint. Excepturi.
          </md-card-content>
        </md-card-expand-content>
      </md-card-expand>
    </md-card>

    <md-card>
      <md-card-media>
        <img src="static/images/girl.png" alt="People">
      </md-card-media>

      <md-card-header>
        <div class="md-title">Title goes here</div>
        <div class="md-subhead">Subtitle here</div>
      </md-card-header>

      <md-card-expand>
        <md-card-actions md-alignment="right">
          <md-card-expand-trigger>
            <md-button>Learn more</md-button>
          </md-card-expand-trigger>
        </md-card-actions>

        <md-card-expand-content>
          <md-card-content>
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Optio itaque ea, nostrum odio.
            Dolores, sed accusantium quasi non, voluptas eius illo quas,
            saepe voluptate pariatur in deleniti minus sint. Excepturi.
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
    questionnaires: null,
  }),
  methods: {
    getQuestionnaires() {
      const path = `${apiUrl}/questionnaires`;
      headers.Authorization = localStorage.getItem('access_token');
      axios.post(path, this.payload, { headers })
        .then((response) => {
          this.questionnaires = JSON.parse(JSON.stringify(response.data));
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
