import React, { Component } from "react";
import {Redirect} from "react-router-dom"
import store from "store2";
import v1 from "../api/v1";


class LoginPage extends Component {
    constructor(props){
        super(props)

        v1.page.set({ title: 'Login' })
        
        this.state = {
            username: '',
            password: '',
            loading: false,
            loggedIn: false,
        }


        this.login = () => {
            if(this.state.username.length < 3) {
                alert('Username should be at least 3 characters long')
                return
            }

            if(this.state.password.length < 4) {
                alert('Password should be atleast 4 characters long')
                return
            }
            this.setState({loading: true})
            v1.security.login({user: this.state.username, password: this.state.password}).then(response => {
                if(response.error) {
                    alert({
                        title: 'Failed to login',
                        text: response.error,
                        icon: 'warning',
                        timer: 10000,
                    })
                    this.setState({loading: false})
                } else if(response._id) {
                    alert({
                        text: 'You have logged in as '+ response.fullname,
                        timer: 1500,
                        toast: true
                    })
                    store.clearAll()
                    if(this.props.reload)
                        this.props.reload(response)
                    this.props.onLogin ? this.props.onLogin(response) : this.setState({loggedIn: true})
                }
            }).catch(_ => {
                this.setState({loading: false})
            })
        }
        v1.page.track()

    }
    render(){
        return this.state.loggedIn ? <Redirect to="/creator/dashboard" /> : (
            
            <>
                <div class="box">
                    <div class="box-header">
                        <h4 class="box-title">Login to your account.</h4>
                    </div>
                </div>
                <div class="row box-body">

                    <div class="col-12 form-group">
                        <label>* Username or email</label>
                        <div class="input-group mb-4">
                            <div class="input-group-prepend">
                                <span class="input-group-text bg-white px-4 border-md border-right-0">
                                    <i class="fa fa-user text-muted"></i>
                                </span>
                            </div>
                            <input id="email" type="text" name="email" value={this.state.username} onChange={evt => this.setState({username: evt.target.value})} placeholder="Email address or username" class="form-control bg-white border-left-0 border-md" />
                        </div>
                    </div>

                    <div class="col-12 form-group">
                        <label>* Password</label>
                        <div class="input-group  mb-4">
                            <div class="input-group-prepend">
                                <span class="input-group-text bg-white px-4 border-md border-right-0">
                                    <i class="fa fa-lock text-muted"></i>
                                </span>
                            </div>
                            <input id="password" type="password"  value={this.state.password} onChange={evt => this.setState({password: evt.target.value})} name="password" placeholder="Password" class="form-control bg-white border-left-0 border-md" />
                        </div>
                    </div>

                    <div class="form-group col-lg-12 mx-auto mb-0">
                        { this.state.loading ? <p>Logging you in...</p> :
                        <button type="submit" onClick={this.login}  class="btn btn-info btn-block py-2">
                            <span class="font-weight-bold text-white">Login to your account</span>
                        </button>
                        }
                    </div>
                </div>
            </>
        )
    }
}

export default LoginPage;