<template>
  <div>
    <h1>Buyer: {{ buyerInfo.Buyer[0].id }}</h1>
    <h2> Name : {{ buyerInfo.Buyer[0].name }} </h2>
    <p> Age : {{ buyerInfo.Buyer[0].age }} </p>
    <div v-if="buyerInfo.BuyersWithSameIP.length > 0">
      <h2>Buyers with the same ip</h2>
      <div v-for="(buyer, index) of buyerInfo.BuyersWithSameIP" :key="index" class="pill-element pa-4 text-center light-green darken-2 text-no-wrap rounded-xl">
        <h3>{{ buyer.name }}</h3>
        <h4>Id: {{ buyer.id }} </h4>
        <p>Age: {{ buyer.age }}</p>
      </div>
    </div>
    <h2>Recommended Products</h2>
    <div v-for="(product, index) of buyerInfo.RecommendedProducts" :key="index" class="pill-element pa-4 text-center light-green darken-2 text-no-wrap rounded-xl">
      <h3>{{ product.name }}</h3>
      <h4>Id: {{ product.id }}</h4>
      <p>Bought {{ product.count }} times by people who also bought the same products</p>
    </div>
  </div>
</template>

<script>

export default {

  async asyncData ({ $axios, params }) {
    const buyerInfo = await $axios.$get(`http://127.0.0.1:5000/buyers/${params.id}`)
    return { buyerInfo }
  }
}
</script>

<style scoped>
  .pill-element{
    margin: 10px;
  }
</style>
