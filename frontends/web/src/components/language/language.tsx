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

import { Component, h, RenderableProps } from 'preact';
import { /* i18nEditorActive, */ extraLanguages } from '../../i18n/i18n';

import { translate as translateDecorator, TranslateProps } from '../../decorators/translate';
import { Dialog } from '../dialog/dialog';
import * as style from './language.module.css';

type TActiveLanguageCodes = 'bg' | 'de' | 'en' | 'es'
    | 'fr' | 'hi' | 'it' | 'ja' | 'ms' | 'pt'
    | 'ru' | 'sl' | 'tr' | 'zh';

type TLanguage = {
    code: TActiveLanguageCodes;
    display: string;
};

type TLanguagesList = TLanguage[];

interface State {
    selectedIndex: number;
    activeDialog: boolean;
    languages: TLanguagesList;
}

interface LanguageSwitchProps {
    languages?: TLanguagesList;
    i18n?: any;
}

type Props = TranslateProps & LanguageSwitchProps;

const defaultLanguages = [
    { code: 'bg', display: 'България' },
    { code: 'de', display: 'Deutsch' },
    { code: 'en', display: 'English' },
    { code: 'es', display: 'Español' },
    { code: 'fa', display: 'فارسی' },
    { code: 'fr', display: 'Français' },
    { code: 'he', display: 'עברית' },
    { code: 'hi', display: 'हिन्दी ' },
    { code: 'it', display: 'Italiano' },
    { code: 'ja', display: '日本語' },
    { code: 'ms', display: 'Bahasa Melayu' },
    { code: 'pt', display: 'Português' },
    { code: 'ru', display: 'Русский' },
    { code: 'sl', display: 'Slovenščina' },
    { code: 'tr', display: 'Türkçe' },
    { code: 'zh', display: '中文' },
] as TLanguagesList;

class LanguageSwitch extends Component<Props, State> {
    constructor(props) {
        super(props);
        const languages = this.props.languages || defaultLanguages;
        if (extraLanguages) {
            extraLanguages.split(',').forEach(code => {
                languages.push({
                    code,
                    display: code,
                } as TLanguage);
            });
        }
        this.state = {
            selectedIndex: this.getSelectedIndex(languages),
            activeDialog: false,
            languages,
        };
    }

    abort = () => {
        this.setState({ activeDialog: false });
    }

    getSelectedIndex = (languages) => {
        const lang = this.props.i18n.language;
        // Check for exact match first.
        let index = languages.findIndex(({ code }) => code === lang);
        // A locale may contain region and other sub tags.
        // Try with a relaxed match, only the first component.
        if (index === -1 && lang.indexOf('-') > 0) {
            const tag = lang.slice(0, lang.indexOf('-'));
            index = languages.findIndex(({ code }) => code === tag);
        }
        if (index === -1 && lang.indexOf('_') > 0) {
            const tag = lang.slice(0, lang.indexOf('_'));
            index = languages.findIndex(({ code }) => code === tag);
        }
        // Give up. We tried.
        if (index === -1) {
            return 0;
        }
        return index;
    }

    componentWillMount() {

        /* if (i18nEditorActive) {
         *     // Get languages from backend instead when translating,
         *     // as new languages won't be shown otherwise.
         *     this.context.
         *         i18n.
         *         services.
         *         backendConnector.
         *         backend.
         *         getLanguages((err, data) => {
         *             if (err) {
         *                 alert(err);
         *                 return;
         *             }
         *             const languages = Object.entries(data).map(([key, value]) => {
         *                 return {
         *                     code: key,
         *                     display: value.nativeName,
         *                 };
         *             });
         *             this.setState({ languages, selectedIndex: this.getSelectedIndex(languages) });
         *         });
         * } */
    }

    changeLanguage = ({ target }) => {
        const langCode = target.dataset.code;
        const index = parseInt(target.dataset.index, 10);
        this.setState({
            selectedIndex: index,
            activeDialog: false,
        });
        this.context.i18n.changeLanguage(langCode);
    }

    render({
        t,
    }: RenderableProps<Props>, {
        selectedIndex,
        activeDialog,
        languages,
    }: State) {
        if (languages.length === 1) {
            return null;
        }
        return (
            <div>
                <button
                    title="Select Language"
                    className={style.link}
                    onClick={() => this.setState({ activeDialog: true })}>
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round">
                        <circle cx="12" cy="12" r="10" />
                        <line x1="2" y1="12" x2="22" y2="12" />
                        <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                    </svg>
                    {languages[selectedIndex].code === 'en' ? 'Other languages' : 'English'}
                </button>
                {
                    activeDialog && (
                        <Dialog small slim title={t('language.title')} onClose={this.abort}>
                            {
                                languages.map((language, i) => {
                                    const selected = selectedIndex === i;
                                    return (
                                        <button
                                            key={language.code}
                                            class={[style.language, selected ? style.selected : ''].join(' ')}
                                            onClick={this.changeLanguage}
                                            data-index={i}
                                            data-code={language.code}>
                                            {language.display}
                                            {
                                                selected && (
                                                    <svg
                                                        className={style.checked}
                                                        xmlns="http://www.w3.org/2000/svg"
                                                        width="24"
                                                        height="24"
                                                        viewBox="0 0 24 24"
                                                        fill="none"
                                                        stroke="currentColor"
                                                        strokeWidth="2"
                                                        strokeLinecap="round"
                                                        strokeLinejoin="round">
                                                        <polyline points="20 6 9 17 4 12" />
                                                    </svg>
                                                )
                                            }
                                        </button>
                                    );
                                })
                            }
                        </Dialog>
                    )
                }
            </div>
        );
    }
}

const TranslatedLanguageSwitcher = translateDecorator<LanguageSwitchProps>()(LanguageSwitch);

export { TranslatedLanguageSwitcher as LanguageSwitch };
