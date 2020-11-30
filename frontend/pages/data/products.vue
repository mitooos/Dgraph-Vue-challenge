<template>
  <div>
    <v-row>
      <v-menu
        v-model="fromDateMenu"
        :close-on-content-click="false"
        :nudge-right="40"
        lazy
        transition="scale-transition"
        offset-y
        full-width
        max-width="290px"
        min-width="290px"
      >
        <template v-slot:activator="{ on }">
          <v-text-field
            label="Date"
            readonly
            :value="fromDateDisp"
            v-on="on"
          />
        </template>
        <v-date-picker
          v-model="fromDateVal"
          locale="en-in"
          no-title
          @input="fromDateMenu = false"
        />
      </v-menu>
    </v-row>
     <v-alert v-if="!fromDateVal" class="red">
      Date is required
    </v-alert>
    <v-file-input
      v-model="file"
      small-chips
      truncate-length="15"
    />
    <v-alert v-if="!file" class="red">
      File is required
    </v-alert>
    <v-text-field
      v-if="loading"
      color="success"
      loading
      disabled
    />
    <v-btn @click="submit">
      Submit
    </v-btn>
  </div>
</template>

<script>
export default {
  data () {
    return {
      fromDateMenu: false,
      fromDateVal: null,
      file: null,
      loading: false
    }
  },
  computed: {
    fromDateDisp () {
      return this.fromDateVal
    }
  },
  methods: {
    submit () {
      if (!this.file || !this.fromDateVal) {
        return
      }
      this.loading = true
      const formData = new FormData()
      formData.append('file', this.file)
      formData.append('date', this.fromDateVal)
      this.$axios.$post('http://localhost:5000/products', formData).then((resp) => {
        alert('Done loading data')
        this.loading = false
      }).catch(err => alert(err))
    }
  }
}
</script>
