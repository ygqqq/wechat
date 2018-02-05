import Vue from 'vue'
import Vuex from 'vuex'

import { Toast } from 'mint-ui'

Vue.use(Vuex)

const state = {
  username: '',
  ygq:'',
  userFriends: []
  // ws 注册WebSockets
  // userFriends 好友列表
}

const mutations = {
  getCookie(state) {
    var name= 'username'
    var arr = document.cookie.match(new RegExp("(^| )"+name+"=([^;]*)(;|$)"));
    if(arr != null){
      state.username = arr[2]
    }else{
      return null;
    }
  },
  agreeFriend (state,otherName) {
    MessageBox.confirm(otherName+'想加你好友，同意吗？').then(() => {
      this.$store.state.ws.send(
        JSON.stringify({
          src: this.username,
          dst: otherName,
          messagetype: AgreeAdd
        }
      ))
    },() => {
      this.$store.state.ws.send(
        JSON.stringify({
          src: this.username,
          dst: otherName,
          messagetype: DisAgreeAdd
        }
      ))
    });
  },
  setState (state,msg) {
    state.userMessage.find( item => item.UserName === msg.src ).Status = msg.messagetype%2
    Toast({
      message: msg.src+ ( msg.messagetype%2 ===1? '上线了':'下线了'),
      position: 'bottom',
      duration: 3000
    })
  },
  reception (state,msg) {
     console.log(state,msg)
  }
}

export default new Vuex.Store({
  state,
  mutations
})