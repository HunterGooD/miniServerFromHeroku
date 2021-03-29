let app = new Vue({
    el: "#app",
    vuetify: new Vuetify(),
    data: {
        dialog: false,
        items: ['Foo', 'Bar', 'Fizz', 'Buzz'],
        links: [
            '',
            '',
            '',
            '',
        ],
        view: "all",
        agents: [{
            ID: 0,
            Login: "asd",
            Storages: [{
                ID: 0,
                NameStorage: "",
                Address: "",
                Autos: [{
                    NameAuto: "",
                    Photos: [],
                }]
            }]
        }]
    },
    computed: {
        Agents() {
            let v = this;
            return v.agents.filter(
                (agent) => {
                    let ret = true;
                    switch (v.view) {
                        case "all":
                            ret = true;
                            break
                        case "onlyLoaded":
                            let newStorages = agent.Storages.filter(
                                (s) => {
                                    let a = s.Autos.filter(
                                        (a) => {
                                            return a.Photos.length == 0 ? false : true
                                        }
                                    );
                                    return a.length == 0 ? false : true
                                }
                            );
                            ret = newStorages.length == 0 ? false : true;
                            break;
                        case "dontLoad":
                            let newStorages1 = agent.Storages.filter(
                                (s) => {
                                    let a = s.Autos.filter(
                                        (a) => {
                                            return a.Photos.length == 0 ? true : false;
                                        }
                                    );
                                    return a.length == 0 ? true : false;
                                }
                            );
                            ret = newStorages1 == 0 ? true : false;
                            break;
                    }
                    return ret;
                }
            );
        },
    },
    methods: {
        loadAgents() {
            let v = this;
            let xhr = new XMLHttpRequest();
            xhr.onload = function () {
                var res = JSON.parse(xhr.response);
                v.agents = res
            };
            xhr.open('GET', "/agents", true);
            xhr.send();
        }
    },
    mounted() {
        this.loadAgents()
    }
});