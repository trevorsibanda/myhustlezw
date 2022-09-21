import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'
import CreatorTopNav from "../components/CreatorNav";
import ErrorBoundary from "../components/ErrorBoundary";
import TopNav from "../components/TopNav";

import PublicCreatorAccountPage from "./PublicCreatorAccountPage";
import PublicCreatorContentPage from "./PublicCreatorContentPage";
import PublicBuyMeACoffee from "./PublicBuyMeACoffee";
import store from "store2";
import PublicCreatorContentList from "./PublicCreatorContentList";
import PublicCreatorSupporters from "./PublicCreatorSupporters";
import v1 from "../api/v1";

class CreatorPage extends Component{
    constructor(props) {
        super(props)
        let { username } = this.props.match.params
        v1.page.set({title: 'Loading...'})
        
        let state = {
            loading: true,
            not_found: false,
            creator: v1.defaults.creator(username),
            featured: v1.defaults.featured,
            headline: v1.defaults.headline,
            campaigns: v1.defaults.campaigns,
            supporters: v1.defaults.supporters,
        }

        if(window.pregenerated_creator && window.pregenerated_creator.creator && window.pregenerated_creator.creator.username === username){
            state = { ...this.state, ...window.pregenerated_creator }
            state.loading = false
        }

            
        this.state = state 
        
        //choose page load strategy
        if (this.state.loading) {
            v1.public.getCreator(username).then(res => {
                if (res.error) {
                    this.setState({
                        loading: false,
                        error: res.error,
                        not_found: true,
                    })
                    alert('Failed to load with error '+ res.error)
                }
                v1.page.set({ title: res.creator.fullname+ ' @' + res.creator.username })
                this.setState(res)
                this.setState({loading: false})
            }).catch(err => {
                this.setState({not_found: true, loading: false})
                //alert('Failed to load with error ' + err.error)
                return
            })    
        }
    }

    render() {
        return (
            <div class="cover-your-area bg-ghostwhite" >
                <ErrorBoundary>
                    <TopNav creator={this.state.creator} user={this.props.user} />
                    <CreatorTopNav creator={this.state.creator} />
                    <Switch>
                        <Route exact path="/@:username" >
                            <PublicCreatorAccountPage
                                creator={this.state.creator}
                                loading={this.state.loading}
                                not_found={this.state.not_found}
                                supporters={this.state.supporters}
                                subscription={this.state.subscription}
                                user={this.props.user}
                                campaigns={this.state.campaigns}
                                featured={this.state.featured}
                                stream_url={this.state.stream_url} />
                        </Route>
                        
                        <Route exact path="/@:username/buymeacoffee" render={props => <PublicBuyMeACoffee math={props.match} user={this.props.user} />} />
                        <Route exact path="/@:username/content" render={props => <PublicCreatorContentList match={props.match} user={this.props.user} creator={this.state.creator} campaigns={this.state.campaigns} />} />
                        <Route exact path="/@:username/supporters" render={props => <PublicCreatorSupporters match={props.match} user={this.props.user} creator={this.state.creator} supporters={this.state.supporters} />} />
                        <Route exact path="/@:username/:content_id" render={
                            (props) => <PublicCreatorContentPage match={props.match} user={this.props.user} creator={this.state.creator} campaigns={this.state.campaigns} />} 
                        />
                    </Switch>
                </ErrorBoundary>
            </div>
        )
    }
}

export default CreatorPage;