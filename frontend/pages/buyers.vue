<template>
  <div>
    <h1 class="text-center">Buyers</h1>
    <div v-for="(buyer, index) of buyers" :key="index" class="buyer-element pa-4 text-center light-green darken-2 text-no-wrap rounded-xl">
      <h2>{{ buyer.name }}</h2>
      <h3>Id: {{ buyer.id }} </h3>
      <p>Age: {{ buyer.age }}</p>
    </div>
  </div>
</template>

<script>
export default {
  async asyncData ({ $axios }) {
    const buyers = await $axios.$get('http://127.0.0.1:5000/buyers')
    return { buyers }
  }
}
</script>

<style scoped>
  .buyer-element{
    margin: 10px;
  }
</style>
