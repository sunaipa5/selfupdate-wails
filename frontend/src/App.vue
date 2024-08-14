<template>
<div>
  <h1>Checking update</h1>
  <h2>Version 0.0.3</h2>
  <button @click="checkUpdate()">CHECK</button>
</div>
</template>

<script>
import { CheckUpdate } from '../wailsjs/go/main/App.js'
export default{
methods:{
  checkUpdate(){
    CheckUpdate()
  }
}
}
</script>

<style>
body{
background-color: #090909;
color: #999;
}
</style>