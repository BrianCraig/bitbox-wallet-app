/**
 * Copyright 2018 Shift Devices AG
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

import { Component, h, RenderableProps } from 'preact';
import { share } from '../../decorators/share';
import { Store } from '../../decorators/store';
import { translate, TranslateProps } from '../../decorators/translate';
import { setConfig } from '../../utils/config';
import { apiGet } from '../../utils/request';
import A from '../anchor/anchor';
import * as style from './guide.module.css';

export interface SharedProps {
    shown: boolean;
    activeSidebar: boolean;
    sidebarStatus: string;
    guideExists: boolean;
}

export const store = new Store<SharedProps>({
    shown: false,
    activeSidebar: false,
    sidebarStatus: '',
    guideExists: false,
});

apiGet('config').then(({ frontend }) => {
    if (frontend && frontend.guideShown != null) { // eslint-disable-line eqeqeq
        store.setState({ shown: frontend.guideShown });
    } else {
        store.setState({ shown: true });
    }
});

function setGuideShown(shown: boolean) {
    store.setState({ shown });
    setConfig({ frontend: { guideShown: shown } });
}

export function toggle() {
    setGuideShown(!store.state.shown);
}

export function show() {
    setGuideShown(true);
}

export function hide() {
    setGuideShown(false);
}

type Props = SharedProps & TranslateProps;

class Guide extends Component<Props> {
    public componentDidMount() {
        store.setState({ guideExists: true });
    }

    public componentWillUnmount() {
        store.setState({ guideExists: false });
    }

    public render(
        { shown, t, children }: RenderableProps<Props>,
    ) {
        return (
            <div className={style.wrapper}>
                <div className={[style.overlay, shown && style.show].join(' ')} onClick={toggle}></div>
                <div className={[style.guide, shown && style.show].join(' ')}>
                    <div className={[style.header, 'flex flex-row flex-between flex-items-center'].join(' ')}>
                        <h2>{t('guide.title')}</h2>
                        <a href="#" className={style.close} onClick={toggle}>
                            {t('guide.toggle.close')}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round">
                                <line x1="18" y1="6" x2="6" y2="18"></line>
                                <line x1="6" y1="6" x2="18" y2="18"></line>
                            </svg>
                        </a>
                    </div>
                    <div className={style.content}>
                        {children}
                        <div className={style.entry}>
                            {t('guide.appendix.text')} <A href={t('guide.appendix.href')}>{t('guide.appendix.link')}</A>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

const HOC = translate()(share<SharedProps, TranslateProps>(store)(Guide));
export { HOC as Guide };
