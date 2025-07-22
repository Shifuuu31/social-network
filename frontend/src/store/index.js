import { createStore } from "vuex";
import chat from "@/store/modules/chat.js"
import notifications from "@/store/modules/notifications.js"
import actions from "@/store/actions.js"

export default createStore({
  modules: {
    chat,
    notifications
  },
  //------------------------------------- state is like a variables, which hold a values.
  state: {
    id: "", // id is currently logged in user ID
    wsConn: {}, // changed from null to {} for logout fix. (vic) 

    profileInfo: {},
    myFollowers: null,

    posts: {
      allposts: [],
      myposts: [],
      groupPosts: [],
    },

    users: {
      allusers: [],
    },

    groups: {
      allGroups: [],
      userGroups: [],
    },


    dataLoaded: {
      userGroups: false
    }

  },
  //------------------------------------ getters is a way for check state values.
  getters: {
    getId(state) {
      return state.id;
    },
    allPosts(state) {
      return state.posts.allposts;
    },
    myPosts(state) {
      return state.posts.myposts;
    },
    groupPosts(state) {
      return state.posts.groupPosts;
    },
    userInfo(state) {
      return state.profileInfo;
    },
    allUsers(state) {
      return state.users.allusers;
    },
    allGroups(state) {
      return state.groups.allGroups;
    },
    followers(state){
      return state.myFollowers
    },
    filterUsers: (state) => (searchquery) => {
      if (searchquery === "") {
        return [];
      }
      let arr = [];
      state.users.allusers.filter((user) => {
        if (user.nickname.toLowerCase().includes(searchquery.toLowerCase())) {
          arr.push(user);
        }
      });
      return arr;
    },

    filterGroups: (state) => (searchquery) => {
      if (searchquery === "") {
        return [];
      }
      let arr = [];
      state.groups.allGroups.filter((group) => {
        if (group.name.toLowerCase().includes(searchquery.toLowerCase())) {
          arr.push(group);
        }
      });
      return arr;
    },

    getMyFollowersNames({ myFollowers }) {
      if (myFollowers === null) {
        return null
      }

      return myFollowers.map((follower) => {
        if (follower.nickname) {
          return follower.nickname
        } else {
          return follower.firstName + follower.lastName
        }
      })
    },

    getMyFollowerIDs({ myFollowers }) {
      if (Array.isArray(myFollowers) && myFollowers.length > 0) {
        return myFollowers.map((follower) => follower.id)
      }
    },

  },
  //-------------------------------------- mutations is a way for change state.
  mutations: {
    updatePosts(state, posts) {
      state.posts.allposts = posts;
    },
    updateMyPosts(state, myposts) {
      state.posts.myposts = myposts;
    },
    updateProfileInfo(state, userinfo) {
      state.profileInfo = userinfo;
    },
    updateAllUsers(state, users) {
      state.users.allusers = users;
    },
    updateAllGroups(state, groups) {
      state.groups.allGroups = groups;
    },

    updateMyFollowers(state, followers) {
      state.myFollowers = followers;
    },

    updateMyUserID(state, id) {
      state.id = id;
    },
    updateGroupPosts(state, posts) {
      state.posts.groupPosts = posts;
    },
    updateWebSocketConn(state, wsConn) {
      state.wsConn = wsConn
    },

    updateUserGroups(state, userGroups) {
      state.groups.userGroups = userGroups
    },


    updateDataLoaded(state, data) {
      state.dataLoaded[data] = true;
    }


  },
  //------------------------------------------Actions
  actions: actions
});
