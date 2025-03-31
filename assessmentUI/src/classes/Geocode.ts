import { type Reactive } from 'vue';
import type { Device } from '../types/types';

type AddressData = {
  address: Address;
  class: 'place' | 'highway' | 'building' | 'landuse';
};

type Address = {
  house_number: string;
  road: string;
  county: string;
  city?: string;
  town?: string;
  state: string;
  postcode: string;
  'ISO3166-2-lvl4': string;
};

export type Place = {
  lat: number;
  lng: number;
  formattedAddress: string;
  lastUpdated: number;
};

type Places = Reactive<{ [deviceId: string]: Place }>;

const streetTypeAbbreviations: { [key: string]: string } = {
  Avenue: 'Ave',
  Boulevard: 'Blvd',
  Circle: 'Cir',
  Court: 'Ct',
  Drive: 'Dr',
  Lane: 'Ln',
  Place: 'Pl',
  Road: 'Rd',
  Street: 'St',
  Way: 'Way',
};

export class Geocode {
  deviceIdCache: string[];
  places: Places;
  reverseUrl: URL;
  userAgent: { 'User-Agent': string };
  lastTimeFetched: number;
  apiCallDelay: number;
  constructor() {
    this.deviceIdCache = [];
    this.places = {};
    this.reverseUrl = new URL('https://nominatim.openstreetmap.org/reverse');
    this.userAgent = { 'User-Agent': 'one-step-gps/1.0 (adnan@ashihabi.me)' };
    this.lastTimeFetched = Date.now();
    this.apiCallDelay = 0;
  }

  async getLocation(device: Device): Promise<void> {
    // prevent geocode api spamming
    if (!areNSecondsApart(Date.now(), this.lastTimeFetched, 2000)) {
      this.lastTimeFetched = Date.now();
      this.apiCallDelay += 2000;
      await new Promise<void>((resolve) =>
        setTimeout(() => {
          this.apiCallDelay -= 2000;
          resolve();
        }, this.apiCallDelay),
      );
    }

    const { deviceId, latestDevicePoint, online } = device;
    const { lat, lng, deviceState } = latestDevicePoint;
    const { driveStatus } = deviceState;
    // lower precision (to around 11 meters) to not exhaust API rate limit
    const roundedLat = parseFloat(lat.toFixed(4));
    const roundedLng = parseFloat(lng.toFixed(4));

    const indexInCache = this.deviceIdCache.indexOf(deviceId);
    // if not in cache
    if (indexInCache === -1) {
      // add to cache if not moving
      if (!online || driveStatus === 'idle' || driveStatus === 'off') {
        this.deviceIdCache.push(deviceId);
      }

      // if coordinates are same as previously stored coordinates:
      // do not fetch new address
      const currentDeviceLocation = this.places[deviceId];

      if (
        currentDeviceLocation &&
        this._isSameLocation(roundedLat, roundedLng, currentDeviceLocation)
      ) {
        return;
      }

      this.lastTimeFetched = Date.now();
      // fetch address based on coordinates
      const location = await this._reverseGeocode(roundedLat, roundedLng);
      if (location) {
        const formattedAddress = this._formatAddress(location);
        this.places[deviceId] = {
          lat: roundedLat,
          lng: roundedLng,
          lastUpdated: Date.now(),
          formattedAddress,
        };
      }
    } else if (online && driveStatus === 'driving') {
      // remove from cache if active again
      this.deviceIdCache.slice(indexInCache, 1);
    }
  }

  async _reverseGeocode(lat: number, lng: number): Promise<AddressData | null> {
    this.reverseUrl.search = new URLSearchParams({
      format: 'json',
      lat: String(lat),
      lon: String(lng),
      addressdetails: '1',
    }).toString();

    try {
      const response = await fetch(this.reverseUrl, {
        headers: this.userAgent,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error(`Error geocoding ${lat}, ${lng}:`, error);
      return null;
    }
  }

  _formatAddress(addressData: AddressData): string {
    const address = addressData.address;
    const addressClass = addressData.class;
    const city = address.city ?? address.town ?? address.county;
    if (['place', 'building', 'landuse'].includes(addressClass)) {
      return `${address.house_number || ''} ${address.road ? this._formatRoad(address.road) : ''}, ${city || ''}, ${address['ISO3166-2-lvl4'].split('-')[1]} ${address.postcode}`.trim();
    } else {
      return `${address.road ? this._formatRoad(address.road) : ''}, ${city || ''}, ${address['ISO3166-2-lvl4'].split('-')[1]} ${address.postcode}`.trim();
    }
  }

  _formatRoad(road: string): string {
    const [streetName, streetType] = road.split(' ');
    const newStreetType = streetTypeAbbreviations[streetType] ?? streetType;
    return `${streetName} ${newStreetType}`;
  }

  _isSameLocation(
    lat: number,
    lng: number,
    currentDeviceLocation: Place,
  ): boolean {
    const moveDecimal = 10_000;

    return (
      Math.floor(currentDeviceLocation.lat * moveDecimal) ===
        Math.floor(lat * moveDecimal) &&
      Math.floor(currentDeviceLocation.lng * moveDecimal) ===
        Math.floor(lng * moveDecimal)
    );
  }

  static splitAtFirstComma(str: string) {
    const index = str.indexOf(',');
    if (index === -1) {
      return [str]; // If no comma, return the original string in an array
    }

    return [str.slice(0, index), str.slice(index + 1)];
  }
}

function areNSecondsApart(
  timestamp1: number,
  timestamp2: number,
  n: number,
): boolean {
  const difference = Math.abs(timestamp1 - timestamp2);
  return difference >= n;
}
