export function getFriends(_this){
  //const _this = this
  axios.get('/api/user/friends/'+ _this.username, {})
  .then(function (response) {
    if (response.data.success) {
      if (response.data.msg !== 'null') {
        let friendsArr = JSON.parse(response.data.msg)
        var userMessage = {
          username: friendsArr[0].Friends[0],
          friends: []
        }
        for (var i = 0; i < friendsArr.length; i++) {
          let friends = {
            UserName: friendsArr[i].UserName,
            NickName: friendsArr[i].NickName,
            Id: friendsArr[i].Id_,
            Status: friendsArr[i].Status,
            CreateAt: friendsArr[i].CreateAt,
          }
          userMessage.friends.push(friends)
        }
        //保存到vuex
        _this.$store.state.userMessage = userMessage
      }
    } else {
      console.log('false')
    }
  })
  .catch(function (error) {
    alert(error);
  })
}