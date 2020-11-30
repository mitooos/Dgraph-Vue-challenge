<template>
  <div>
    <v-toolbar>
      <v-toolbar-title class="title">
        <NuxtLink to="/">
          Home
        </NuxtLink>
      </v-toolbar-title>
      <v-toolbar-items>
        <v-btn text to="/buyers">
          Buyers
        </v-btn>
        <v-menu offset-y>
          <template v-slot:activator="{ on, attrs }">
            <v-btn text dark to="#" v-bind="attrs" v-on="on">
              Sync data
            </v-btn>
          </template>
          <v-list>
            <v-btn text to="/data/products">
              Products
            </v-btn>
            <v-btn text to="/data/buyers">
              Buyers
            </v-btn>
            <v-btn text to="/data/transactions">
              Transactions
            </v-btn>
          </v-list>
        </v-menu>
      </v-toolbar-items>
    </v-toolbar>
  </div>
</template>

<style scoped>
  .title{
    margin-right: 10px;
  }
</style>
