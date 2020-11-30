<template>
  <div>
    <h1 class="text-center">
      Buyers
    </h1>
    <div v-for="(buyer, index) of buyers" :key="index" class="buyer-element pa-4 text-center light-green darken-2 text-no-wrap rounded-xl">
      <NuxtLink :to="'/buyers/'+buyer.id">
        <h2>{{ buyer.name }}</h2>
        <h3>Id: {{ buyer.id }} </h3>
        <p>Age: {{ buyer.age }}</p>
      </NuxtLink>
    </div>
    <div class="text-center">
      <v-pagination
        v-model="page"
        :length="6"
      />
    </div>
  </div>
</template>

<script>
export default {
  async asyncData ({ $axios }) {
    const buyers = await $axios.$get('http://127.0.0.1:5000/buyers?first=10&offset=0')
    return { buyers }
  },
  data () {
    return {
      page: 1
    }
  },
  watch: {
    // eslint-disable-next-line object-shorthand
    page: function (val) {
      this.page = val
      this.fetchData()
    }
  },
  methods: {
    fetchData () {
      this.$axios.$get(`http://127.0.0.1:5000/buyers?first=10&offset=${this.page * 10}`).then((resp) => {
        this.buyers = resp
        window.scrollTo(0, 0)
      })
    }
  }
}
</script>

<style scoped>
  .buyer-element{
    margin: 10px;
  }
  a{
    color: white !important;
  }
</style>
