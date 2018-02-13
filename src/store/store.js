import Vue from 'vue'
import Vuex from 'vuex'

import { Toast } from 'mint-ui'
import { MessageBox } from 'mint-ui'

Vue.use(Vuex)

const ErrorMsg = 0  //错误消息
const OnlineRemind	= 1	//上线提醒
const OfflineRemind   = 2 //下线提醒 
const AddFriendReq	= 3 //添加好友请求
const AgreeAdd		= 4 //同意好友请求
const DisAgreeAdd 	= 5 //拒绝好友请求
const ChatMsg			= 6 //普通聊天消息
const NormalMsg		= 10 //普通通知消息

const state = {
  username: '',
  ygq:'',
  userFriends: [],
  nowMessage: []
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
     state.ws.send(
        JSON.stringify({
          src: state.username,
          dst: otherName,
          messagetype: AgreeAdd
        }
      ))
    },() => {
      state.ws.send(
        JSON.stringify({
          src: state.username,
          dst: otherName,
          messagetype: DisAgreeAdd
        }
      ))
    });
  },
  setState (state,msg) {
    state.userFriends.find( item => item.UserName === msg.src ).Status = msg.messagetype%2
    Toast({
      message: msg.src+ ( msg.messagetype%2 ===1? '上线了':'下线了'),
      position: 'bottom',
      duration: 3000
    })
  },
  reception (state,msg) {
     console.log(state,msg)
  },
  chatMessages (state,msg) {
    state.nowMessage.unshift(msg)
  }
}

export default new Vuex.Store({
  state,
  mutations
})