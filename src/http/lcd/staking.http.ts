import {HttpService, Injectable} from '@nestjs/common';
import {cfg} from '../../config/config';
import {Logger} from '../../logger';
import {
    AddressBalancesLcdDto,
    IconUriLcdDto,
    StakingValidatorDelegationLcdDto,
    StakingValidatorLcdDto, StakingValidatorParametersLcdDto,
    StakingValidatorSlashLcdDto, StakingValUnBondingDelLcdDto,
    DelegatorsDelegationLcdDto,
    DelegatorsUndelegationLcdDto,
    BondedTokensLcdDto,
    IDelegationLcd,
    TokensStakingLcdToken
} from "../../dto/http.dto";

@Injectable()

export class StakingHttp {
    async queryValidatorListFromLcd(status: string, pageNum: number, pageSize: number) {
        let validatorLcdUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/validators?status=${status}&pagination.limit=${pageSize}`;
        try {
            let stakingValidatorData: any = await new HttpService().get(validatorLcdUri).toPromise().then(result => result.data)
            if (stakingValidatorData && stakingValidatorData.validators) {
                return StakingValidatorLcdDto.bundleData(stakingValidatorData.validators);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${validatorLcdUri}`, e)
            throw new Error(e);
        }
    }

    async queryValidatorFormSlashing(address_ica: string) {
        // const slashValidatorUri = `${cfg.serverCfg.lcdAddr}/slashing/validators/${validatorPubkey}/signing_info`
        //todo:hangtaishan 测试网api
        const slashValidatorUri = `${cfg.serverCfg.lcdAddr}/cosmos/slashing/v1beta1/signing_infos/${address_ica}`
        try {
            const stakingSlashValidatorData: any = await new HttpService().get(slashValidatorUri).toPromise().then(result => result.data)
            if (stakingSlashValidatorData && stakingSlashValidatorData.val_signing_info) {
                return new StakingValidatorSlashLcdDto(stakingSlashValidatorData.val_signing_info);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            if (e && e.response && e.response.data && e.response.data.code == 2) {
                Logger.warn(`api-error from ${slashValidatorUri}`, e.response.data)
            } else {
                Logger.warn(`api-error from ${slashValidatorUri}`, e)
            }
        }
    }

    async queryValidatorIcon(valIdentity) {
        const getIconUri = `${cfg.serverCfg.iconUri}?fields=pictures&key_suffix=${valIdentity || ''}`
        try {
            const valIconData: any = await new HttpService().get(getIconUri).toPromise().then(result => result.data)
            if (valIconData) {
                return new IconUriLcdDto(valIconData)
            } else {
                Logger.warn('api-error:', 'there is no result of validators from getIconUri');
            }

        } catch (e) {
            Logger.warn(`api-error from ${getIconUri}`, e)
        }
    }

    async queryParametersFromSlashing() {
        const parameterUri = `${cfg.serverCfg.lcdAddr}/cosmos/slashing/v1beta1/params`
        try {
            const parameterData: any = await new HttpService().get(parameterUri).toPromise().then(result => result.data)
            if (parameterData && parameterData.params) {
                return new StakingValidatorParametersLcdDto(parameterData.params)
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }

        } catch (e) {
            Logger.warn(`api-error from ${parameterUri}`, e)
        }
    }

    async queryValidatorDelegationsFromLcd(address, pageNum=1, pageSize=1000, useCount=false) {
        // const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/validators/${address}/delegations`;
        let offset = (Number(pageNum) - 1) * Number(pageSize);
        const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/validators/${address}/delegations?pagination.offset=${offset}&pagination.limit=${pageSize}&pagination.count_total=${useCount}`;
        try {
            let validatorDelegationsData: any = await new HttpService().get(getValidatorDelegationsUri).toPromise().then(result => result.data)
            if (validatorDelegationsData && validatorDelegationsData.delegation_responses) {
                let data = { total: validatorDelegationsData.pagination && Number(validatorDelegationsData.pagination.total), result: validatorDelegationsData.delegation_responses };
                return new StakingValidatorDelegationLcdDto(data);
            } else {
                Logger.warn('api-error:', 'there is no result of validator delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getValidatorDelegationsUri}`, e)
        }
    }

    async queryValidatordelegatorNumFromLcd(address) {
        // const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/validators/${address}/delegations`
        const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/validators/${address}/delegations?pagination.limit=1&pagination.count_total=true`
        try {
            let validatorDelegationsData: any = await new HttpService().get(getValidatorDelegationsUri).toPromise().then(result => result.data)
            if (validatorDelegationsData && validatorDelegationsData.pagination && validatorDelegationsData.pagination.total) {
                return Number(validatorDelegationsData.pagination.total);
            } else {
                Logger.warn('api-error:', 'there is no result of validator delegatorNum from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getValidatorDelegationsUri}`, e)
        }
    }

    async queryValidatorSelfBondFromLcd(address_iva,address_iaa) {
        const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/validators/${address_iva}/delegations/${address_iaa}`
        try {
            let validatorDelegationsData: any = await new HttpService().get(getValidatorDelegationsUri).toPromise().then(result => result.data)
            if (validatorDelegationsData && validatorDelegationsData.delegation_response) {
                return new IDelegationLcd(validatorDelegationsData.delegation_response);
            } else {
                Logger.warn('api-error:', 'there is no result of validator delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getValidatorDelegationsUri}`, e)
        }
    }



    async queryValidatorUnBondingDelegations(address, pageNum=1, pageSize=1000, useCount=false) {
        // const getValidatorUnBondingDelUri = `${cfg.serverCfg.lcdAddr}/staking/validators/${address}/unbonding_delegations`
        let offset = (Number(pageNum) - 1) * Number(pageSize);
        const getValidatorUnBondingDelUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/validators/${address}/unbonding_delegations?pagination.offset=${offset}&pagination.limit=${pageSize}&pagination.count_total=${useCount}`
        try {
            let validatorUnBondingDelegationsData: any = await new HttpService().get(getValidatorUnBondingDelUri).toPromise().then(result => result.data)
            if (validatorUnBondingDelegationsData && validatorUnBondingDelegationsData.unbonding_responses) {
                let data = { total: validatorUnBondingDelegationsData.pagination && Number(validatorUnBondingDelegationsData.pagination.total), result: validatorUnBondingDelegationsData.unbonding_responses };
                return new StakingValUnBondingDelLcdDto(data);
            } else {
                Logger.warn('api-error:', 'there is no result of validator unBonding delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getValidatorUnBondingDelUri}`, e)
        }
    }

    async queryBalanceByAddress(address) {
        const getBalancesUri = `${cfg.serverCfg.lcdAddr}/cosmos/bank/v1beta1/balances/${address}`
        try {
            let addressBalancesData: any = await new HttpService().get(getBalancesUri).toPromise().then(result => result.data)
            if (addressBalancesData && addressBalancesData.balances) {
                return AddressBalancesLcdDto.bundleData(addressBalancesData.balances);
            } else {
                Logger.warn('api-error:', 'there is no result of validator unBonding delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getBalancesUri}`, e)
        }
    }

    async queryDelegatorsDelegationsFromLcd(address, pageNum=1, pageSize=1000, useCount=false) {
        // const getDelegatorsDelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/delegators/${address}/delegations`;
        let offset = (Number(pageNum) - 1) * Number(pageSize);
        const getDelegatorsDelegationsUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/delegations/${address}?pagination.offset=${offset}&pagination.limit=${pageSize}&pagination.count_total=${useCount}`
        try {
            const delegatorsDelegationsData: any = await new HttpService().get(getDelegatorsDelegationsUri).toPromise().then(result => result.data)
            if (delegatorsDelegationsData && delegatorsDelegationsData.delegation_responses) {
                return new DelegatorsDelegationLcdDto(delegatorsDelegationsData);
            } else {
                Logger.warn('api-error:', 'there is no result of delegators delegations from lcd');
            }
        } catch (e) {
            if (e && e.response && e.response.data && e.response.data.code == 2) {
                Logger.warn(`api-error from ${getDelegatorsDelegationsUri}`, e.response.data)
            } else {
                Logger.warn(`api-error from ${getDelegatorsDelegationsUri}`, e)
            }
        }
    }

    async queryDelegatorsUndelegationsFromLcd(address, pageNum=1, pageSize=1000, useCount=false) {
        // const getDelegatorsUndelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/delegators/${address}/unbonding_delegations`
        let offset = (Number(pageNum) - 1) * Number(pageSize);
        const getDelegatorsUndelegationsUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/delegators/${address}/unbonding_delegations?pagination.offset=${offset}&pagination.limit=${pageSize}&pagination.count_total=${useCount}`
        try {
            const delegatorsUndelegationsData: any = await new HttpService().get(getDelegatorsUndelegationsUri).toPromise().then(result => result.data)
            if (delegatorsUndelegationsData && delegatorsUndelegationsData.unbonding_responses) {
                return new DelegatorsUndelegationLcdDto(delegatorsUndelegationsData);
            } else {
                Logger.warn('api-error:', 'there is no result of delegators delegations from lcd');
            }
        } catch (e) {
            if (e && e.response && e.response.data && e.response.data.code == 2) {
                Logger.warn(`api-error from ${getDelegatorsUndelegationsUri}`, e.response.data)
            } else {
                Logger.warn(`api-error from ${getDelegatorsUndelegationsUri}`, e)
            }
        }
    }

    static async getBondedTokens () {
        const BondedTokensUrl = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/pool`
        try {
            let BondedTokens: any = await new HttpService().get(BondedTokensUrl).toPromise().then(result => result.data)
            if (BondedTokens && BondedTokens.pool) {
                return new BondedTokensLcdDto(BondedTokens.pool);
            } else {
                Logger.warn('api-error:', 'there is no result of bonded_tokens from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${BondedTokensUrl}`, e)
        }
    }

    async getStakingTokens() {
        const stakingTokensUrl = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/params`
        try {
            let stakingTokensData: any = await new HttpService().get(stakingTokensUrl).toPromise().then(result => result.data)
            if (stakingTokensData && stakingTokensData.params) {
                return new TokensStakingLcdToken(stakingTokensData.params);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${stakingTokensUrl}`, e)
        }
    }
}
