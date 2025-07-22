import { createRouter, createWebHistory } from "vue-router";
import Auth from "../components/Auth.vue";
import store from "@/store";
// import SignIn from '../views/SignInView.vue'
// import RegisterView from '../views/RegisterView.vue'


const routes = [
  {
    path: "/",
    name: "auth",
    component: Auth,

  },
  {
    path: "/sign-in",
    name: "sign-in",
    component: () => import("../views/SignInView.vue"),
  },
  {
    path: "/reg",
    name: "register",
    component: () => import("../views/RegisterView.vue"),
  },
  {
    path: "/main",
    name: "mainpage",
    components: {
      default: () => import("../views/MainView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue")
    }
  },
  {
    path: "/profile/:id",
    name: "Profile",
    components: {
      default: () => import("../views/ProfileView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue")
    }
  },
  {
    path: "/group/:id",
    name: "Group",
    components: {
      default: () => import("../views/GroupView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue")
    }
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});


router.beforeEach(async (to, from) => {
  const isAuthenticated = await store.dispatch("isLoggedIn");

  // if user is not authenticated redirect back to sign in page BUT
  // only if the page user wants to go is not sign-in or register
  if (!isAuthenticated && to.name !== "sign-in" && to.name !== "register") {
    return { name: "sign-in" }
  }
})

export default router;
