<template>
  <v-app id="inspire">
    <v-dialog
      v-model="dialog"
      fullscreen
      hide-overlay
      transition="dialog-bottom-transition"
    >
      <v-card>
        <v-toolbar dark color="primary">
          <v-btn icon dark @click="dialog = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
          <v-toolbar-title>Просмотр</v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>
        <v-list three-line subheader />
        <v-img
          contain
          height="95vh"
          :src="dialog_img"
          class="white--text align-end"
          gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
        >
        </v-img>
      </v-card>
    </v-dialog>

    <v-main class="grey lighten-3">
      <v-container>
        <v-row>
          <v-col cols="2">
            <v-sheet rounded="lg">
              <v-list color="transparent">
                <v-list-item
                  v-for="n in storages"
                  :key="n.id"
                  link
                  @click="openStorage(n.id)"
                >
                  <v-list-item-content>
                    <v-list-item-title>
                      Склад: {{ n.name_storage }}
                    </v-list-item-title>
                    <v-list-item-content class="text--secondary text-caption">
                      Улица: {{ n.address }}
                    </v-list-item-content>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </v-sheet>
          </v-col>

          <v-col>
            <v-sheet min-height="95vh" rounded="lg" style="padding: 1% 3%">
              <v-select
                :items="items"
                label="Выберите объект"
                @change="select"
              ></v-select>
              <v-spacer></v-spacer>

              <v-col v-for="e in active_photos" :key="e.id">
                <v-card elevation="8">
                  <v-img
                    :src="getURL(e.path)"
                    aspect-ratio="16/9"
                    height="400px"
                    contain
                    @click="showDialog(e.path)"
                  >
                    <!-- <v-card-title :v-text=""></v-card-title> -->
                  </v-img>
                  <v-divider></v-divider>
                  <v-card-subtitle>
                    Пользователь: {{ e.user.fio }} <v-divider></v-divider> Координаты:
                    Longitude - {{ e.longitude }}; Latitude - {{ e.latitude }} <v-divider></v-divider>
                    Дата: {{ new Date(e.created_at).toISOString().slice(0,10) }}
                  </v-card-subtitle>
                </v-card>
              </v-col>
            </v-sheet>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: "Home",
  data: () => ({
    dialog_img: "https://picsum.photos/200/300",
    dialog: false,
    url: "/api",
    items: [],
    user: "",
    active_photos: [
    ],
    active_storage: [
      {
        name_object: "",
        photos: [
          {
            id: 0,
            created_at: "",
            longitude: "",
            latitude: "",
            path: "",
          },
        ],
      },
    ],
    storages: [],
    agents: [
      {
        id: 0,
        fio: "asd",
        storages: [
          {
            id: 0,
            name_storage: "",
            Address: "",
            objects: [
              {
                name_object: "",
                photos: [
                  {
                    id: 0,
                    created_at: "",
                    longitude: "",
                    latitude: "",
                    path: "",
                  },
                ],
              },
            ],
          },
        ],
      },
    ],
  }),
  computed: {},
  methods: {
    getURL(url) {
      return `${this.$router.app.baseURL}/photo/${
        url[0] === "/" ? url.replace("/", "") : url
      }`;
    },
    showDialog(url) {
      this.dialog_img = url;
      this.dialog = true;
    },
    select(e) {
      let v = this;
      v.active_storage.forEach(async (element) => {
        if (element.name_object == e) {
          v.active_photos = element.photos;
        }
      });
    },
    openStorage(id) {
      let v = this;
      let xhr = new XMLHttpRequest();
      v.agents.forEach((e) => {
        let uF = e.fio;
        let f = false;
        e.storages.forEach((el) => {
          if (el.id == id) {
            f = true;
            return;
          }
        });
        if (f) {
          v.user = uF;
          return;
        }
      });
      xhr.onload = function () {
        var res = JSON.parse(xhr.response);
        v.active_storage = res;
        v.items = [];
        res.forEach((e) => {
          v.items.push(e.name_object);
        });
      };
      xhr.open("GET", `${v.url}/storage/${id}`, true);
      xhr.send();
    },
    loadAgents() {
      let v = this;
      let xhr = new XMLHttpRequest();
      xhr.onload = function () {
        var res = JSON.parse(xhr.response);
        v.agents = res;
        v.agents.forEach((e) => {
          v.storages.push(...e.storages);
        });
      };
      xhr.open("GET", `${v.url}/storages`, true);
      xhr.send();
    },
  },
  mounted() {
    this.loadAgents();
  },
};
</script>
