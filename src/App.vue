<template>
  <div id="app">
    <router-view/>
    <button v-if="Show_flexNav"><router-link to="/login">登录</router-link></button>
    <button v-if="Show_flexNav"><router-link to="/register">注册</router-link></button>
    <ul class='flexNav' v-if="Show_flexNav">
      <router-link to="/message" tag="li">
      <span class="icon-bubble2 icon-common"></span>
      <span>消息</span>
      </router-link>
      <!-- <router-link to="/addres" tag="li">通讯录</router-link> -->
      <router-link to="/home" tag="li">
      <span class="icon-user icon-common"></span>
      <span>我</span>
      </router-link>
    </ul>
  </div>
</template>

<script>
import config from "../config/local.config"
import { mapState,mapMutations,mapActions } from "vuex"
import { getFriends } from './content/script/getFriends'

export default {
  name: 'App',
  data () {
    return {
    }
  },
  methods: {
    ...mapMutations([
      'getCookie', // 将 `this.increment()` 映射为 `this.$store.commit('increment')`
    ])
  },
  computed: {
    ...mapState({
      username: 'username',
    }),
    Show_flexNav () {
      let status = this.$route.fullPath
      if( status == "/login" || status == "/register" || this.$route.name == "chat"){
        return false
      }else{
        return true
      }
    }
  },
  created () {
    this.getCookie()
    let path = this.$route.path
    let _this = this
     
    if (this.username) {
      if ( path!== '/login' && path!== '/register') {
        if (!this.$store.state.ws) {
          this.$store.state.ws = new WebSocket(config.wsUrl+'?a='+this.username) //注册WebSockets
        }
      }
    } else{
      this.$router.push({ name: 'login'});
    }
  }
}
</script>

<style lang="scss"  type="text/css">
@import url('./content/style/reset.css');
@import url('./icon/style.css');

.router-link-active{
  color:#26a2ff;
}
*{
	box-sizing: border-box;
}
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
  background: #f6f6f6;
}

.flexNav{
  position: fixed;
  bottom: 0;
  left: 0;
  height: 50px;
  line-height: 50px;
  display: flex;
  width: 100%;
  background: #fff;

  li{
    flex: 1;
  }

  /* todo 公用不变形边框线 */flex: 1;
    font-size: 0;
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
    font-size: 0;

    span{
      display: block;
      font-size: 12px;
      margin: -18px 0 0 0;
      line-height: 12px;

      &.icon-common{
        font-size: 22px;
        line-height: 22px;
        margin:4px 0 0 0;
      }
    }
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
