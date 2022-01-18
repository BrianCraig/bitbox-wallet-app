/**
 * Copyright 2018 Shift Devices AG
 * Copyright 2021 Shift Crypto AG
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { FunctionComponent, useCallback, useEffect, useMemo } from 'react';
import { AppRouter } from './routes/router'
import { syncAccountsList } from './api/accountsync';
import { syncDeviceList } from './api/devicessync';
import { ConnectedApp } from './connected';
import { Alert } from './components/alert/Alert';
import { Aopp } from './components/aopp/aopp';
import { Banner } from './components/banner/banner';
import { Confirm } from './components/confirm/Confirm';
import { store as panelStore } from './components/guide/guide';
import { MobileDataWarning } from './components/mobiledatawarning';
import { Sidebar, toggleSidebar } from './components/sidebar/sidebar';
import { Update } from './components/update/update';
import i18nEditorActive from './i18n/i18n';
import { apiPost } from './utils/request';
import { apiWebsocket } from './utils/websocket';
import { RouterWatcher } from './utils/route';
import { useSubscribe } from './hooks/api';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

/**
 * Gets fired when the route changes.
 */
const handleRoute = () => {
    if (panelStore.state.activeSidebar) {
        toggleSidebar();
    }
}

const emptyAccounts = [];
const emptyDevices = {};

export const App: FunctionComponent = () => {
    const { t } = useTranslation();
    const route = useNavigate();

    // Show Notifications
    const ws = useSubscribe(apiWebsocket);
    useEffect(() => {
        if(ws === undefined){
            return;
        }
        const { type, data, meta } = (ws as any);
        switch (type) {
        case 'backend':
            switch (data) {
            case 'newTxs':
                apiPost('notify-user', {
                    text: t('notification.newTxs', {
                        count: meta.count,
                        accountName: meta.accountName,
                    }),
                });
                break;
            }
            break;
        }
    }, [ws, t]);

    const accounts = useSubscribe(syncAccountsList) || emptyAccounts;
    const devices = useSubscribe(syncDeviceList) || emptyDevices;

    const deviceList = useMemo(() => Object.keys(devices), [devices]);

    const maybeRoute = useCallback(() => {
        const currentURL = window.location.pathname;
        const isIndex = currentURL === '/' || currentURL === '/index.html' || currentURL === '/android_asset/web/index.html';
        const inAccounts = currentURL.startsWith('/account/');

        // QT and Android start their apps in '/index.html' and '/android_asset/web/index.html' respectively
        // This re-routes them to '/' so we have a simpler uri structure
        if (isIndex && currentURL !== '/' && (!accounts || accounts.length === 0)) {
            route('/', {replace: true});
            return;
        }
        // if no accounts are registered on specified views route to /
        if ( accounts.length === 0 && (
            currentURL.startsWith('/account-summary')
            || currentURL.startsWith('/add-account')
            || currentURL.startsWith('/settings/manage-accounts')
            || currentURL.startsWith('/passphrase')
        )) {
            route('/', {replace: true});
            return;
        }
        // if on an account that isnt registered route to /
        if (inAccounts && !accounts.some(account => currentURL.startsWith('/account/' + account.code))) {
            route('/', {replace: true});
            return;
        }
        // if on index page and there is at least 1 account route to /account-summary
        if (isIndex && accounts && accounts.length) {
            route('/account-summary', {replace: true});
            return;
        }
        // if on the /buy/ view and there are no accounts view route to /
        if (accounts.length === 0 && currentURL.startsWith('/buy/')) {
            route('/', {replace: true});
            return;
        }
    }, [accounts, route]);
    
    useEffect(() => {
        // without accounts route to device settings for unlock, pair, create, restore etc.
        if(deviceList.length > 0){
            if (accounts.length) {
                maybeRoute();
                return;
            }
            route(`/device/${deviceList[0]}`, {replace: true});
        } else {
            maybeRoute();
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [devices, accounts]);
 
    // Returns a string representation of the current devices, so it can be used in the `key` property of subcomponents.
    // The prefix is used so different subcomponents can have unique keys to not confuse the renderer.
    const devicesKey = (prefix: string): string => {
        return prefix + ':' + JSON.stringify(devices, deviceList.sort());
    }
 
    const activeAccounts = accounts.filter(acct => acct.active);

    return (
        <ConnectedApp>
            <div className={['app', i18nEditorActive ? 'i18nEditor' : ''].join(' ')}>
                <Sidebar
                    accounts={activeAccounts}
                    deviceIDs={deviceList} />
                <div className="appContent flex flex-column flex-1" style={{minWidth: 0}}>
                    <Update />
                    <Banner msgKey="bitbox01" />
                    <MobileDataWarning />
                    <Aopp />
                    <AppRouter
                        accounts={accounts}
                        activeAccounts={activeAccounts}
                        deviceIDs={deviceList}
                        devices={devices}
                        devicesKey={devicesKey}
                    />
                    <RouterWatcher onChange={handleRoute} />
                </div>
                <Alert />
                <Confirm />
            </div>
        </ConnectedApp>
    );
}
