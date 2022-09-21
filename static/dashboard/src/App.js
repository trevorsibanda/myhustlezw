import React, { Component, Suspense } from "react"
import { Container } from "react-bootstrap";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  useHistory,
} from "react-router-dom";

import './App.css';

import SideNav from './components/SideNav'
import TopNav from './components/TopNav'

import v1 from "./api/v1";


import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import swal from 'sweetalert';

import Preloader from "./components/PreLoader";

import BottomNav from "./components/BottomNavBar";
import ScrollToTop from "./components/ScrollToTop";
import ErrorBoundary from "./components/ErrorBoundary";
import CreatorPage from "./pages/CreatorPage";
import store from "store2";
import ReactGA from 'react-ga';
import FloatingActionButton from "./components/FloatingActionButton";
import { Redirect } from "react-router-dom/cjs/react-router-dom.min";

const CreatorDashboard = React.lazy(() => import("./pages/CreatorDashboard"));
const AuthPage = React.lazy(() => import("./pages/AuthPage"));


class App extends Component  {

  constructor(props) {
    super(props)

    this.floatingRoutes = [
      {icon: "fa-home", route: "/creator/dashboard", targetRoute: "/creator/dashboard/timeline", targetIcon: "fa-newspaper"},
      {icon: "fa-heart", route: "/creator/supporters", targetRoute: "/creator/supporters/subscriptions", targetIcon: "fa-key"},
      {icon: "fa-credit-card", route: "/creator/wallet", targetRoute: "/creator/wallet/mypayments", targetIcon: "fa-money"},
      {icon: "fa-plus", route: "/creator/content/new", targetRoute: "/creator/content", targetIcon: "fa-th-list"},
    ]

    
    v1.payments.syncCurrenyRate().then(() => {
    }).catch(_ => { })

    

    window.alert = function (options) {
      swal(options)
    }

    window.supportsHLS = _ => {
      var video = document.createElement('video');
      return Boolean(video.canPlayType('application/vnd.apple.mpegURL') || video.canPlayType('audio/mpegurl'))
    }

    this.logout = function(evt) {
      if(window.location.pathname.startsWith("/auth")){
        evt.preventDefault()
        return
      }
      v1.user.logout()
      store.clearAll()
      swal({
        toast: true,
        title: "Logged out",
        text: "You have been logged out of your account.",
        icon: "warning",
        timer: 1500,
      })
      setTimeout(_ => {
        useHistory().push("/auth/login")
      }, 3000)
      evt.preventDefault()
    }

    window.logout = this.logout

    //sync currency every 12 hours
    setInterval(v1.payments.syncCurrenyRate, 60000 * 60 * 12)



    TimeAgo.addDefaultLocale(en)
    

    this.state = {
      loading: true,
      loggedIn: true,
    }


    this.load = () => {
      //load all preloads
      if(window.pregenerated && window.pregenerated.user && window.pregenerated.user._id !== undefined){
        this.setState({ user: window.pregenerated.user, loading: false })
        return v1.user.current(false).then(user => {
            this.setState({ user: user, loading: false })
          }).catch(err => {
            if (err && err.error === "Not authenticated") {
              window.logout()
              this.setState({ loading: false, loggedIn: false })
            }
          })
      } else {
        v1.user.current(true).then(user => {
          this.setState({ user: user, loading: false })
        }).catch(_ => {
          return v1.user.current(false).then(user => {
            this.setState({ user: user, loading: false })
          }).catch(err => {
            if (err && err.error === "Not authenticated") {
              window.logout()
              this.setState({ loading: false, loggedIn: false })
            }
          })
        })
      }
      
    }
    this.load()
    ReactGA.initialize(v1.config.gaTrackingCode)
    v1.page.set({
      title: 'Dashboard',
    })

  }


  render() {
    return this.state.loading && ! window.location.pathname.startsWith('/auth') ? <Preloader /> : 
    <div class="container">
      <Router>
        <ScrollToTop>
          
          <Switch>
            <Route path="/auth" >
                
                <ErrorBoundary >
                  <Suspense fallback={<Preloader />} >
                    <AuthPage user={this.state.user} reload={user => { this.setState({ user, loading: false }); this.load(); }} />
                  </Suspense>
                </ErrorBoundary>
            </Route>
            <Route path="/@:username" render={(props) => <CreatorPage user={this.state.user} match={props.match} />} />
            <Route path="/_r/:creator/:resource" render={(props) =>  <Redirect to={"/" + props.match.params.creator + "/" + props.match.params.resource} /> } />
            <Route path="/" >
              <TopNav user={this.state.user} />
              <div class="padding-top-80 padding-bottom-110" >
                <div class="row padding-top-10">
                  <div class="col-md-4 col-lg-3">
                    <SideNav user={this.state.user} />
                  </div>
                  <div class="col-lg-9 col-md-8">
                      <Suspense fallback={<p>Loading dashboard...</p>}>
                        <CreatorDashboard user={this.state.user} />
                      </Suspense>
                      <BottomNav user={this.state.user} />
                      
                  </div>
                </div>
              </div>

            </Route>
          </Switch>

          </ScrollToTop>
          <FloatingActionButton routes={this.floatingRoutes} user={this.state.user} />
        </Router>
    </div>
    
  }
}

export default App;
