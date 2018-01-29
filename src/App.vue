<template>
  <div id="app">
    <router-view/>
    <ul class='flexNav' v-if="Show_flexNav">
      <router-link to="/message" tag="li">消息</router-link>
      <router-link to="/addres" tag="li">通讯录</router-link>
      <router-link to="/home" tag="li">我</router-link>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'App',
  data () {
    return {
    }
  },
  computed: {
    Show_flexNav () {
      let status = this.$route.fullPath
      if( status == "/login" || status == "/register"){
        return false
      }else{
        return true
      }
    }          
  }
}
</script>

<style lang="scss"  type="text/css">
@import url('./content/style/reset.css');

body, html{
  width: 100%;
  height: 100%;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 100%;
  height: 100%;
}

.flexNav{
  position: fixed;
  bottom: 0;
  left: 0;
  height: 50px;
  line-height: 50px;
  display: flex;
  width: 100%;

  li{
    flex: 1;
  }

  /* todo 公用不变形边框线 */
  &::before{
    height: 1PX;
    content: '';
    width: 100%;
    border-bottom: 1PX solid #f0f0f0;
    position: absolute;
    top: -1PX;
    right: 0;
    transform: scaleY(0.5);
    -webkit-transform: scaleY(0.5);
    z-index: 10;
  }
}

.flexNav{
  position: fixed;
  bottom: 0;
  left: 0;
  height: 50px;
  line-height: 50px;
  display: flex;
  width: 100%;

  li{
    flex: 1;
  }

  /* todo 公用不变形边框线 */
  &::before{
    height: 1PX;
    content: '';
    width: 100%;
    border-bottom: 1PX solid #f0f0f0;
    position: absolute;
    top: -1PX;
    right: 0;
    transform: scaleY(0.5);
    -webkit-transform: scaleY(0.5);
    z-index: 10;
  }
}
</style>
