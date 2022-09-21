import { Redirect, Route, Switch } from 'react-router-dom';
import ErrorBoundary from '../components/ErrorBoundary';
import React, { Suspense } from 'react';


const Dashboard = React.lazy(() => import("./Dashboard"))
const Supporters = React.lazy(() => import("./Supporters"))
const Campaigns = React.lazy(() => import("./Campaigns"))
const Wallet = React.lazy(() => import("./Wallet"))
const UserSettings = React.lazy(() => import("./UserSettings"));
const EditCampaign = React.lazy(() => import("./EditCampaign"));
const FanSubscriptions = React.lazy(() => import("./FanSubscriptions"));
const FanPayments = React.lazy(() => import("./FanPayments"));
const FanUpgradeAccount = React.lazy(() => import("./FanUpgradeAccount"));
const ViewSupporter = React.lazy(() => import("./ViewSupporter"));
const FanContentFeed = React.lazy(() => import("./FanContentFeed"));
const ViewPayment = React.lazy(() => import("./ViewPayment"));
const CreateNewCampaign = React.lazy(() => import("./CreateCampaign"));

/*
{  this.state.user && !this.state.user.phoneVerified ?
                        <ActionCenter user={this.state.user} action='phone' /> :
                        <>}
*/
function CreatorDashboard(props) {
    return (
        <ErrorBoundary>
            <Suspense fallback={<p>Loading...</p>} >
            <Switch>
                <Route exact path="/creator/wallet/payment/:id" render={(props1) => <ViewPayment match={props1.match} user={props.user} />} />
                            
                <Route exact path="/creator/content/:id/order/:order_id"></Route>
                <Route exact path="/creator/content/new">
                    <CreateNewCampaign user={props.user} />
                </Route>
                <Route exact path="/creator/content/:id" render={(props1) => <EditCampaign match={props1.match} user={props.user} />} />
                
                <Route exact path="/creator/content">
                    <Campaigns user={props.user} />
                </Route>
                
                <Route exact path="/creator/wallet/mypayments">
                    <FanPayments user={props.user} />
                </Route>
                <Route exact path="/creator/wallet">
                    <Wallet user={props.user} />
                </Route>
                <Route path="/creator/settings"><UserSettings user={props.user} /></Route>
                <Route exact path="/creator/supporters/subscriptions">
                    <FanSubscriptions user={props.user} />
                </Route>
                <Route exact path="/creator/supporters/:id" render={(props1) => <ViewSupporter match={props1.match} user={props.user} />} />
                
                <Route exact path="/creator/supporters">
                    <Supporters user={props.user} />
                </Route>
                
                <Route exact path="/creator/help">
                    
                </Route>
                <Route exact path="/creator/verify-identity">
                    <FanUpgradeAccount user={props.user} />
                </Route>
            
                <Route exact path="/creator/dashboard/timeline">
                    <FanContentFeed user={props.user} />
                </Route>
                <Route exact path="/creator/dashboard">
                    <Dashboard user={props.user} />
                </Route>
                <Route path="/creator" >
                    <Redirect to="/creator/dashboard" />
                </Route>

            </Switch>  
            </Suspense>            
        </ErrorBoundary>
    )
}

export default CreatorDashboard;