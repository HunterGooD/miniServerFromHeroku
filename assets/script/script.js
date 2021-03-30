let app = new Vue({
    el: "#app",
    vuetify: new Vuetify(),
    data: {
        dialog_img: "https://picsum.photos/200/300",
        dialog: false,
        items: [],
        user: "",
        active_photos: [],
        active_storage: [{
            name_object: "",
            photos: [{
                id: 0,
                created_at: "",
                longitude: "",
                latitude: "",
                path: "",
            }],
        }],
        storages: [],
        agents: [{
            id: 0,
            fio: "asd",
            storages: [{
                id: 0,
                name_storage: "",
                Address: "",
                objects: [{
                    name_object: "",
                    photos: [{
                        id: 0,
                        created_at: "",
                        longitude: "",
                        latitude: "",
                        path: "",
                    }],
                }]
            }]
        }]
    },
    computed: {
        url(hash) {
            return "https://defsgthjyhtgrkj.herokuapp.com/photo/" + hash;
        }
    },
    methods: {
        showDialog(url) {
            this.dialog_img = url;
            this.dialog = true;
        },
        select(e) {
            let v = this;
            v.active_storage.forEach(async element => {
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
                res.forEach(e => {
                    v.items.push(e.name_object)
                });
            };
            xhr.open('GET', "/api/storage/" + id, true);
            xhr.send();
        },
        loadAgents() {
            let v = this;
            let xhr = new XMLHttpRequest();
            xhr.onload = function () {
                var res = JSON.parse(xhr.response);
                v.agents = res;
                v.agents.forEach(e => {
                    v.storages.push(...e.storages);
                });
            };
            xhr.open('GET', "/api/storages", true);
            xhr.send();
        }
    },
    mounted() {
        this.loadAgents()
    }
});